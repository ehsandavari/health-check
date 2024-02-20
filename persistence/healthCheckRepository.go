package persistence

import (
	"github.com/ehsandavari/go-logger"
	"health-check/application/interfaces"
	"health-check/domain/entities"
	"health-check/infrastructure/postgres"
	"health-check/pkg/genericRepository"
	"health-check/pkg/tracer"
)

type sHealthCheckRepository struct {
	iLogger   logger.ILogger
	iTracer   tracer.ITracer
	sPostgres postgres.SPostgres
	genericRepository.IGenericRepository[entities.HealthCheck]
}

func NewHealthCheckRepository(logger logger.ILogger, tracer tracer.ITracer, postgres postgres.SPostgres) interfaces.IHealthCheckRepository {
	return sHealthCheckRepository{
		iLogger:            logger,
		iTracer:            tracer,
		sPostgres:          postgres,
		IGenericRepository: genericRepository.NewGenericRepository[entities.HealthCheck](logger, tracer, postgres),
	}
}
