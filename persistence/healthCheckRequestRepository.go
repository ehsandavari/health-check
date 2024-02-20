package persistence

import (
	"github.com/ehsandavari/go-logger"
	"health-check/application/interfaces"
	"health-check/domain/entities"
	"health-check/infrastructure/postgres"
	"health-check/pkg/genericRepository"
	"health-check/pkg/tracer"
)

type sHealthCheckRequestRepository struct {
	iLogger   logger.ILogger
	iTracer   tracer.ITracer
	sPostgres postgres.SPostgres
	genericRepository.IGenericRepository[entities.HealthCheckRequest]
}

func NewHealthCheckRequestRepository(logger logger.ILogger, tracer tracer.ITracer, postgres postgres.SPostgres) interfaces.IHealthCheckRequestRepository {
	return sHealthCheckRequestRepository{
		iLogger:            logger,
		iTracer:            tracer,
		sPostgres:          postgres,
		IGenericRepository: genericRepository.NewGenericRepository[entities.HealthCheckRequest](logger, tracer, postgres),
	}
}
