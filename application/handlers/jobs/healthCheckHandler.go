package jobs

import (
	"encoding/json"
	"fmt"
	"github.com/ehsandavari/go-context-plus"
	"github.com/ehsandavari/go-logger"
	"health-check/application/interfaces"
	"health-check/domain/entities"
	"health-check/domain/enums"
	"health-check/pkg/genericRepository"
	"health-check/pkg/tracer"
	"net/http"
)

type SHealthCheckJobHandler struct {
	iLogger       logger.ILogger
	iTracer       tracer.ITracer
	iRedis        interfaces.IRedis
	iCron         interfaces.ICron
	iRest         interfaces.IRest
	iNotification interfaces.INotification
	iUnitOfWork   interfaces.IUnitOfWork

	callAddJob           func(ctx *contextplus.Context, healthCheck entities.HealthCheck)
	callSubRedis         func(ctx *contextplus.Context)
	callSendRequest      func(ctx *contextplus.Context, healthCheck entities.HealthCheck)
	callSendNotification func(ctx *contextplus.Context, subject string, msg string)

	healthCheckChannel chan string
}

func newHealthCheckJobHandler(
	iLogger logger.ILogger,
	iTracer tracer.ITracer,
	iRedis interfaces.IRedis,
	iCron interfaces.ICron,
	iRest interfaces.IRest,
	iNotification interfaces.INotification,
	iUnitOfWork interfaces.IUnitOfWork,
) SHealthCheckJobHandler {
	s := SHealthCheckJobHandler{
		iLogger:            iLogger,
		iTracer:            iTracer,
		iRedis:             iRedis,
		iCron:              iCron,
		iRest:              iRest,
		iNotification:      iNotification,
		iUnitOfWork:        iUnitOfWork,
		healthCheckChannel: make(chan string),
	}
	s.callAddJob = s.addJob
	s.callSubRedis = s.subRedis
	s.callSendRequest = s.sendRequest
	s.callSendNotification = s.sendNotification
	return s
}

func (r SHealthCheckJobHandler) Start(ctx *contextplus.Context) error {
	span, ctx := r.iTracer.SpanFromContext(ctx)
	defer span.Finish()

	var healthChecks []entities.HealthCheck

	healthChecks, err := r.iUnitOfWork.HealthCheckRepository().All(ctx, genericRepository.Equal("status", enums.StatusStart))
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("err", err)
		r.iLogger.WithError(err).Error(ctx, "error in get all health checks")

		return err
	}

	for _, healthCheck := range healthChecks {
		r.callAddJob(ctx, healthCheck)
	}

	go r.callSubRedis(ctx)

	return nil
}

func (r SHealthCheckJobHandler) Stop(ctx *contextplus.Context) error {
	span, ctx := r.iTracer.SpanFromContext(ctx)
	defer span.Finish()

	close(r.healthCheckChannel)

	return nil
}

func (r SHealthCheckJobHandler) addJob(ctx *contextplus.Context, healthCheck entities.HealthCheck) {
	span, ctx := r.iTracer.SpanFromContext(ctx)
	defer span.Finish()

	if healthCheck.Status == enums.StatusStop || healthCheck.DeletedAt.Valid {
		r.iCron.RemoveJob(healthCheck.Id)
		return
	}

	if err := r.iCron.AddJob(healthCheck.Id, healthCheck.UpdatedAt, healthCheck.Interval, func() {
		r.callSendRequest(ctx, healthCheck)
	}); err != nil {
		span.SetTag("error", true)
		span.LogKV("err", err)
		r.iLogger.WithError(err).Error(ctx, "error in add job")

		return
	}
}

func (r SHealthCheckJobHandler) subRedis(ctx *contextplus.Context) {
	span, ctx := r.iTracer.SpanFromContext(ctx)
	defer span.Finish()

	go r.iRedis.Subscribe(ctx, "healthCheck", r.healthCheckChannel)

	var healthCheck entities.HealthCheck

	for c := range r.healthCheckChannel {
		if err := json.Unmarshal([]byte(c), &healthCheck); err != nil {
			span.SetTag("error", true)
			span.LogKV("err", err)
			r.iLogger.WithError(err).WithString("data", c).Error(ctx, "error in json unmarshal to health check entity")

			continue
		}

		r.callAddJob(ctx, healthCheck)
	}
}

func (r SHealthCheckJobHandler) sendRequest(ctx *contextplus.Context, healthCheck entities.HealthCheck) {
	span, ctx := r.iTracer.SpanFromContext(ctx)
	defer span.Finish()

	statusCode, header, body, err := r.iRest.Execute(ctx, healthCheck.Method, healthCheck.Url, healthCheck.Headers.Data(), healthCheck.Body.Data())
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("err", err)
		r.iLogger.WithError(err).Error(ctx, "error in execute rest request")

		return
	}

	healthCheckRequest := entities.NewHealthCheckRequest(
		healthCheck.Id,
		header,
		body,
		statusCode,
	)
	if err = r.iUnitOfWork.HealthCheckRequestRepository().Create(ctx, &healthCheckRequest); err != nil {
		span.SetTag("error", true)
		span.LogKV("err", err)
		r.iLogger.WithError(err).WithAny("healthCheckRequest", healthCheckRequest).Error(ctx, "error in create health check request")

		return
	}

	if statusCode != http.StatusOK {
		r.callSendNotification(
			ctx,
			fmt.Sprintf("id : %d | url : %s | method : %s", healthCheck.Id, healthCheck.Url, healthCheck.Method),
			fmt.Sprintf("request id : %d | status code : %d | response body : %s", healthCheckRequest.Id, statusCode, body),
		)
	}
}

func (r SHealthCheckJobHandler) sendNotification(ctx *contextplus.Context, subject string, msg string) {
	span, ctx := r.iTracer.SpanFromContext(ctx)
	defer span.Finish()

	if err := r.iNotification.Send(
		ctx,
		subject,
		msg,
	); err != nil {
		span.SetTag("error", true)
		span.LogKV("err", err)
		r.iLogger.WithError(err).WithString("subject", subject).WithString("msg", msg).Error(ctx, "error in send health check notification")

		return
	}
}
