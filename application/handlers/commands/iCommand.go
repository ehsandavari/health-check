package commands

import (
	"github.com/ehsandavari/go-context-plus"
	"health-check/domain/entities"
	"health-check/infrastructure"
	"health-check/persistence"
)

type ICommand[T1 any, T2 any] interface {
	Handle(ctx *contextplus.Context, command T1) (T2, error)
}

type Commands struct {
	HealthCheckCreate ICommand[SHealthCheckCreateCommand, *entities.HealthCheck]
	HealthCheckDelete ICommand[SHealthCheckDeleteCommand, *entities.HealthCheck]
	HealthCheckStatus ICommand[SHealthCheckStatusCommand, *entities.HealthCheck]
}

func NewCommands(infrastructure *infrastructure.Infrastructure, persistence *persistence.Persistence) Commands {
	return Commands{
		HealthCheckCreate: newHealthCheckCreateCommandHandler(infrastructure.ILogger, infrastructure.ITracer, infrastructure.IRedis, persistence.IUnitOfWork),
		HealthCheckDelete: newHealthCheckDeleteCommandHandler(infrastructure.ILogger, infrastructure.ITracer, infrastructure.IRedis, persistence.IUnitOfWork),
		HealthCheckStatus: newHealthCheckStatusCommandHandler(infrastructure.ILogger, infrastructure.ITracer, infrastructure.IRedis, persistence.IUnitOfWork),
	}
}
