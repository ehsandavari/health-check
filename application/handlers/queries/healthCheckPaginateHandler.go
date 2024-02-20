package queries

import (
	"github.com/ehsandavari/go-context-plus"
	"github.com/ehsandavari/go-logger"
	"health-check/application/common"
	"health-check/application/interfaces"
	"health-check/domain/entities"
	"health-check/pkg/tracer"
)

type SHealthCheckPaginateQueryHandler struct {
	iLogger                logger.ILogger
	iTracer                tracer.ITracer
	iHealthCheckRepository interfaces.IHealthCheckRepository
}

func newHealthCheckPaginateQueryHandler(
	iLogger logger.ILogger,
	iTracer tracer.ITracer,
	iHealthCheckRepository interfaces.IHealthCheckRepository,
) SHealthCheckPaginateQueryHandler {
	return SHealthCheckPaginateQueryHandler{
		iLogger:                iLogger,
		iTracer:                iTracer,
		iHealthCheckRepository: iHealthCheckRepository,
	}
}

func (r SHealthCheckPaginateQueryHandler) Handle(ctx *contextplus.Context, query SHealthCheckPaginateQuery) (*common.PaginateResult[entities.HealthCheck], error) {
	span, ctx := r.iTracer.SpanFromContext(ctx)
	defer span.Finish()

	totalRows, healthChecks, err := r.iHealthCheckRepository.Paginate(
		ctx,
		query.paginateQuery,
	)
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("err", err)
		r.iLogger.WithError(err).WithAny("query", query).Error(ctx, "error in paginate health checks")

		return nil, common.ErrorInternalServer
	}

	return common.NewPaginateResult(healthChecks, query.paginateQuery.GetPage(), query.paginateQuery.GetPerPage(), uint64(totalRows)), nil
}
