package jobs

import (
	"github.com/ehsandavari/go-context-plus"
	"health-check/infrastructure"
	"health-check/persistence"
)

type IJob interface {
	Start(ctx *contextplus.Context) error
	Stop(ctx *contextplus.Context) error
}

type Jobs struct {
	HealthCheck IJob
}

func NewJobs(infrastructure *infrastructure.Infrastructure, persistence *persistence.Persistence) Jobs {
	return Jobs{
		HealthCheck: newHealthCheckJobHandler(infrastructure.ILogger, infrastructure.ITracer, infrastructure.IRedis, infrastructure.ICron, infrastructure.IRest, infrastructure.INotification, persistence.IUnitOfWork),
	}
}
