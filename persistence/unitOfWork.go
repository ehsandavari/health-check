package persistence

import (
	"github.com/ehsandavari/go-context-plus"
	"github.com/ehsandavari/go-logger"
	"gorm.io/gorm"
	"health-check/application/interfaces"
	"health-check/infrastructure/postgres"
	"health-check/pkg/tracer"
)

type sUnitOfWork struct {
	logger                        logger.ILogger
	tracer                        tracer.ITracer
	postgres                      postgres.SPostgres
	iHealthCheckRepository        interfaces.IHealthCheckRepository
	iHealthCheckRequestRepository interfaces.IHealthCheckRequestRepository
}

func NewUnitOfWork(
	logger logger.ILogger,
	tracer tracer.ITracer,
	postgres postgres.SPostgres,
	healthCheckRepository interfaces.IHealthCheckRepository,
	healthCheckRequestRepository interfaces.IHealthCheckRequestRepository,
) interfaces.IUnitOfWork {
	return &sUnitOfWork{
		logger:                        logger,
		tracer:                        tracer,
		postgres:                      postgres,
		iHealthCheckRepository:        healthCheckRepository,
		iHealthCheckRequestRepository: healthCheckRequestRepository,
	}
}

func (r sUnitOfWork) HealthCheckRepository() interfaces.IHealthCheckRepository {
	return r.iHealthCheckRepository
}

func (r sUnitOfWork) HealthCheckRequestRepository() interfaces.IHealthCheckRequestRepository {
	return r.iHealthCheckRequestRepository
}

func (r sUnitOfWork) Do(ctx *contextplus.Context, unitOfWorkBlock func(interfaces.IUnitOfWork) error) error {
	span, ctx := r.tracer.SpanFromContext(ctx)
	defer span.Finish()

	if err := r.postgres.Database.Transaction(func(tx *gorm.DB) error {
		r.postgres.Database = tx
		return unitOfWorkBlock(r)
	}); err != nil {
		span.SetTag("error", true)
		span.LogKV("err", err)

		return err
	}

	return nil
}
