package queries

import (
	"github.com/ehsandavari/go-context-plus"
	"health-check/application/common"
	"health-check/domain/entities"
	"health-check/infrastructure"
	"health-check/persistence"
)

type IQuery[T1 any, T2 any] interface {
	Handle(ctx *contextplus.Context, command T1) (T2, error)
}

type Queries struct {
	HealthCheckPaginate IQuery[SHealthCheckPaginateQuery, *common.PaginateResult[entities.HealthCheck]]
}

func NewQueries(infrastructure *infrastructure.Infrastructure, persistence *persistence.Persistence) Queries {
	return Queries{
		HealthCheckPaginate: newHealthCheckPaginateQueryHandler(infrastructure.ILogger, infrastructure.ITracer, persistence.IHealthCheckRepository),
	}
}
