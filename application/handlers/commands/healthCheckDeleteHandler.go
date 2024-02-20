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

type SHealthCheckDeleteCommandHandler struct {
	iLogger     logger.ILogger
	iTracer     tracer.ITracer
	iRedis      interfaces.IRedis
	iUnitOfWork interfaces.IUnitOfWork
}

func newHealthCheckDeleteCommandHandler(
	iLogger logger.ILogger,
	iTracer tracer.ITracer,
	iRedis interfaces.IRedis,
	iUnitOfWork interfaces.IUnitOfWork,
) SHealthCheckDeleteCommandHandler {
	return SHealthCheckDeleteCommandHandler{
		iLogger:     iLogger,
		iTracer:     iTracer,
		iRedis:      iRedis,
		iUnitOfWork: iUnitOfWork,
	}
}

func (r SHealthCheckDeleteCommandHandler) Handle(ctx *contextplus.Context, command SHealthCheckDeleteCommand) (healthCheck *entities.HealthCheck, err error) {
	span, ctx := r.iTracer.SpanFromContext(ctx)
	defer span.Finish()

	if err = r.iUnitOfWork.Do(ctx, func(iUnitOfWork interfaces.IUnitOfWork) error {
		if healthCheck, err = iUnitOfWork.HealthCheckRepository().SingleOrDefault(
			ctx,
			genericRepository.Equal("id", command.id),
		); err != nil {
			span.SetTag("error", true)
			span.LogKV("err", err)
			r.iLogger.WithError(err).WithUint("id", command.id).Error(ctx, "error in find health check")

			return common.ErrorInternalServer
		}

		if healthCheck == nil {
			return common.ErrorNotFound
		}

		if healthCheck, err = iUnitOfWork.HealthCheckRepository().Delete(
			ctx,
			healthCheck,
			genericRepository.Equal("id", command.id),
		); err != nil {
			span.SetTag("error", true)
			span.LogKV("err", err)
			r.iLogger.WithError(err).WithUint("id", command.id).Error(ctx, "error in delete health check")

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
