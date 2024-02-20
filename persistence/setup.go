package persistence

import (
	"health-check/application/interfaces"
	"health-check/infrastructure"
)

type Persistence struct {
	IHealthCheckRepository        interfaces.IHealthCheckRepository
	IHealthCheckRequestRepository interfaces.IHealthCheckRequestRepository
	IUnitOfWork                   interfaces.IUnitOfWork
}

func NewPersistence(infrastructure *infrastructure.Infrastructure) *Persistence {
	healthCheckRepository := NewHealthCheckRepository(infrastructure.ILogger, infrastructure.ITracer, infrastructure.SPostgres)
	healthCheckRequestRepository := NewHealthCheckRequestRepository(infrastructure.ILogger, infrastructure.ITracer, infrastructure.SPostgres)
	return &Persistence{
		IHealthCheckRepository:        healthCheckRepository,
		IHealthCheckRequestRepository: healthCheckRequestRepository,
		IUnitOfWork:                   NewUnitOfWork(infrastructure.ILogger, infrastructure.ITracer, infrastructure.SPostgres, healthCheckRepository, healthCheckRequestRepository),
	}
}
