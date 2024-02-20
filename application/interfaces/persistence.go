package interfaces

import (
	"github.com/ehsandavari/go-context-plus"
	"health-check/domain/entities"
	"health-check/pkg/genericRepository"
)

//go:generate mockgen -destination=./persistence_mock.go -package=interfaces . IHealthCheckRepository,IHealthCheckRequestRepository,IUnitOfWork

type IHealthCheckRepository interface {
	genericRepository.IGenericRepository[entities.HealthCheck]
}

type IHealthCheckRequestRepository interface {
	genericRepository.IGenericRepository[entities.HealthCheckRequest]
}

type IUnitOfWork interface {
	HealthCheckRepository() IHealthCheckRepository
	HealthCheckRequestRepository() IHealthCheckRequestRepository
	Do(*contextplus.Context, func(IUnitOfWork) error) error
}
