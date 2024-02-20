package commands

import (
	"encoding/json"
	"github.com/ehsandavari/go-context-plus"
	"github.com/ehsandavari/go-logger"
	"health-check/application/common"
	"health-check/application/interfaces"
	"health-check/domain/entities"
	"health-check/pkg/genericRepository"
	"health-check/pkg/tracer"
)

type SHealthCheckStatusCommandHandler struct {
	iLogger     logger.ILogger
	iTracer     tracer.ITracer
	iRedis      interfaces.IRedis
	iUnitOfWork interfaces.IUnitOfWork
}

func newHealthCheckStatusCommandHandler(
	iLogger logger.ILogger,
	iTracer tracer.ITracer,
	iRedis interfaces.IRedis,
	iUnitOfWork interfaces.IUnitOfWork,
) SHealthCheckStatusCommandHandler {
	return SHealthCheckStatusCommandHandler{
		iLogger:     iLogger,
		iTracer:     iTracer,
		iRedis:      iRedis,
		iUnitOfWork: iUnitOfWork,
	}
}

func (r SHealthCheckStatusCommandHandler) Handle(ctx *contextplus.Context, command SHealthCheckStatusCommand) (healthCheck *entities.HealthCheck, err error) {
	span, ctx := r.iTracer.SpanFromContext(ctx)
	defer span.Finish()

	if err = r.iUnitOfWork.Do(ctx, func(iUnitOfWork interfaces.IUnitOfWork) error {
		if healthCheck, err = iUnitOfWork.HealthCheckRepository().SingleOrDefault(
			ctx,
			genericRepository.Equal("id", command.id),
		); err != nil {
			span.SetTag("error", true)
			span.LogKV("err", err)
			r.iLogger.WithError(err).WithUint("id", command.id).WithString("status", command.status.String()).Error(ctx, "error in find health check")

			return common.ErrorInternalServer
		}

		if healthCheck == nil {
			return common.ErrorNotFound
		}

		if healthCheck.Status == command.status {
			return common.ErrorBadRequest
		}

		healthCheck.SetStatus(command.status)

		if healthCheck, err = iUnitOfWork.HealthCheckRepository().Update(
			ctx,
			healthCheck,
			genericRepository.Equal("id", command.id),
		); err != nil {
			span.SetTag("error", true)
			span.LogKV("err", err)
			r.iLogger.WithError(err).WithUint("id", command.id).WithString("status", command.status.String()).Error(ctx, "error in update health check status")

			return common.ErrorInternalServer
		}

		var payload []byte
		if payload, err = json.Marshal(healthCheck); err != nil {
			return common.ErrorInternalServer
		}

		if err = r.iRedis.Publish(ctx, "healthCheck", payload); err != nil {
			return common.ErrorInternalServer
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return healthCheck, nil
}
