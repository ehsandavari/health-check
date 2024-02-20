package commands

import (
	"encoding/json"
	"github.com/ehsandavari/go-context-plus"
	"github.com/ehsandavari/go-logger"
	"health-check/application/common"
	"health-check/application/interfaces"
	"health-check/domain/entities"
	"health-check/domain/enums"
	"health-check/pkg/tracer"
)

type SHealthCheckCreateCommandHandler struct {
	iLogger     logger.ILogger
	iTracer     tracer.ITracer
	iRedis      interfaces.IRedis
	iUnitOfWork interfaces.IUnitOfWork
}

func newHealthCheckCreateCommandHandler(
	iLogger logger.ILogger,
	iTracer tracer.ITracer,
	iRedis interfaces.IRedis,
	iUnitOfWork interfaces.IUnitOfWork,
) SHealthCheckCreateCommandHandler {
	return SHealthCheckCreateCommandHandler{
		iLogger:     iLogger,
		iTracer:     iTracer,
		iRedis:      iRedis,
		iUnitOfWork: iUnitOfWork,
	}
}

func (r SHealthCheckCreateCommandHandler) Handle(ctx *contextplus.Context, command SHealthCheckCreateCommand) (*entities.HealthCheck, error) {
	span, ctx := r.iTracer.SpanFromContext(ctx)
	defer span.Finish()

	healthCheck := entities.NewHealthCheck(command.interval, command.url, command.method, command.headers, command.body, enums.StatusStart)
	if err := r.iUnitOfWork.Do(ctx, func(iUnitOfWork interfaces.IUnitOfWork) error {
		if err := iUnitOfWork.HealthCheckRepository().Create(ctx, &healthCheck); err != nil {
			span.SetTag("error", true)
			span.LogKV("err", err)
			r.iLogger.WithError(err).WithAny("command", command).WithAny("healthCheck", healthCheck).Error(ctx, "error in create new health check")

			return common.ErrorInternalServer
		}

		payload, err := json.Marshal(healthCheck)
		if err != nil {
			return common.ErrorInternalServer
		}

		if err = r.iRedis.Publish(ctx, "healthCheck", payload); err != nil {
			return common.ErrorInternalServer
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &healthCheck, nil
}
