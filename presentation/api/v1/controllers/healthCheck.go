package controllers

import (
	"github.com/ehsandavari/go-context-plus"
	"github.com/ehsandavari/go-logger"
	"github.com/gin-gonic/gin"
	"health-check/application"
	"health-check/application/common"
	"health-check/application/handlers/commands"
	"health-check/application/handlers/queries"
	"health-check/domain/entities"
	"health-check/pkg/apiHandler"
	"health-check/pkg/tracer"
	"health-check/presentation/api/v1/dtos"
)

type sHealthCheckController struct {
	apiHandler.SBaseController
	application *application.Application
}

func NewHealthCheckController(application *application.Application, routerGroup *gin.RouterGroup, iLogger logger.ILogger, iTracer tracer.ITracer) {
	healthCheckController := sHealthCheckController{
		SBaseController: apiHandler.NewBaseController(iLogger, iTracer),
		application:     application,
	}

	routerGroup = routerGroup.Group("/health-check")
	{
		routerGroup.POST("/", apiHandler.BaseController[common.PaginateQuery, *common.PaginateResult[entities.HealthCheck]](healthCheckController.list).Handle(healthCheckController.ILogger))
		routerGroup.POST("/create", apiHandler.BaseController[dtos.HealthCheckCreateRequest, *dtos.HealthCheckCreateResponse](healthCheckController.create).Handle(healthCheckController.ILogger))
		routerGroup.PATCH("/:id/:status", apiHandler.BaseController[dtos.HealthCheckStatusRequest, *dtos.HealthCheckStatusResponse](healthCheckController.status).Handle(healthCheckController.ILogger))
		routerGroup.DELETE("/:id", apiHandler.BaseController[dtos.HealthCheckDeleteRequest, *dtos.HealthCheckDeleteResponse](healthCheckController.delete).Handle(healthCheckController.ILogger))
	}
}

// @Tags		health-check
// @Accept		json
// @Produce	json
// @Param		Accept-Language	header		string					true	"header"	Enums(en, fa)
// @Param		params			body		common.PaginateQuery	true	"body"
// @Success	200				{object}	apiHandler.BaseApiResponse[common.PaginateResult[entities.HealthCheck]]
// @Failure	400				{object}	apiHandler.BaseApiResponse[apiHandler.ApiError]
// @Failure	500				{object}	apiHandler.BaseApiResponse[apiHandler.ApiError]
// @Router		/health-check/ [POST]
func (r *sHealthCheckController) list(ctx *contextplus.Context, dto common.PaginateQuery) (*common.PaginateResult[entities.HealthCheck], error) {
	span, ctx := r.ITracer.SpanFromContext(ctx)
	defer span.Finish()

	healthChecks, err := r.application.Queries.HealthCheckPaginate.Handle(ctx, queries.NewHealthCheckPaginateQuery(
		dto,
	))
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("err", err)
		r.ILogger.WithError(err).WithAny("dto", dto).Error(ctx, "error in send mediator update health check status")

		return nil, err
	}

	return healthChecks, nil
}

// @Tags		health-check
// @Accept		json
// @Produce	json
// @Param		Accept-Language	header		string							true	"header"	Enums(en, fa)
// @Param		params			body		dtos.HealthCheckCreateRequest	true	"body"
// @Success	200				{object}	apiHandler.BaseApiResponse[dtos.HealthCheckCreateResponse]
// @Failure	400				{object}	apiHandler.BaseApiResponse[apiHandler.ApiError]
// @Failure	500				{object}	apiHandler.BaseApiResponse[apiHandler.ApiError]
// @Router		/health-check/create [POST]
func (r *sHealthCheckController) create(ctx *contextplus.Context, dto dtos.HealthCheckCreateRequest) (*dtos.HealthCheckCreateResponse, error) {
	span, ctx := r.ITracer.SpanFromContext(ctx)
	defer span.Finish()

	healthCheck, err := r.application.Commands.HealthCheckCreate.Handle(ctx, commands.NewHealthCheckCreateCommand(
		dto.Interval, dto.Url, dto.Method, dto.Headers, dto.Body,
	))
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("err", err)
		r.ILogger.WithError(err).WithAny("dto", dto).Error(ctx, "error in send mediator create health check")

		return nil, err
	}

	return &dtos.HealthCheckCreateResponse{
		Id:        healthCheck.Id,
		Interval:  healthCheck.Interval,
		Url:       healthCheck.Url,
		Method:    healthCheck.Method,
		Headers:   healthCheck.Headers.Data(),
		Body:      healthCheck.Body.Data(),
		Status:    healthCheck.Status,
		CreatedAt: healthCheck.CreatedAt,
	}, nil
}

// @Tags		health-check
// @Accept		json
// @Produce	json
// @Param		Accept-Language	header		string							true	"header"	Enums(en, fa)
// @Param		params			body		dtos.HealthCheckStatusRequest	true	"body"
// @Success	200				{object}	apiHandler.BaseApiResponse[dtos.HealthCheckStatusResponse]
// @Failure	400				{object}	apiHandler.BaseApiResponse[apiHandler.ApiError]
// @Failure	500				{object}	apiHandler.BaseApiResponse[apiHandler.ApiError]
// @Router		/health-check/:id/:status [PATCH]
func (r *sHealthCheckController) status(ctx *contextplus.Context, dto dtos.HealthCheckStatusRequest) (*dtos.HealthCheckStatusResponse, error) {
	span, ctx := r.ITracer.SpanFromContext(ctx)
	defer span.Finish()

	healthCheck, err := r.application.Commands.HealthCheckStatus.Handle(ctx, commands.NewHealthCheckStatusCommand(
		dto.Id, dto.Status,
	))
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("err", err)
		r.ILogger.WithError(err).WithAny("dto", dto).Error(ctx, "error in send mediator update health check status")

		return nil, err
	}

	return &dtos.HealthCheckStatusResponse{
		Id:        healthCheck.Id,
		Status:    healthCheck.Status,
		UpdatedAt: healthCheck.UpdatedAt,
	}, nil
}

// @Tags		health-check
// @Accept		json
// @Produce	json
// @Param		Accept-Language	header		string							true	"header"	Enums(en, fa)
// @Param		params			body		dtos.HealthCheckDeleteRequest	true	"body"
// @Success	200				{object}	apiHandler.BaseApiResponse[dtos.HealthCheckDeleteResponse]
// @Failure	400				{object}	apiHandler.BaseApiResponse[apiHandler.ApiError]
// @Failure	500				{object}	apiHandler.BaseApiResponse[apiHandler.ApiError]
// @Router		/health-check/:id [DELETE]
func (r *sHealthCheckController) delete(ctx *contextplus.Context, dto dtos.HealthCheckDeleteRequest) (*dtos.HealthCheckDeleteResponse, error) {
	span, ctx := r.ITracer.SpanFromContext(ctx)
	defer span.Finish()

	healthCheck, err := r.application.Commands.HealthCheckDelete.Handle(ctx, commands.NewHealthCheckDeleteCommand(
		dto.Id,
	))
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("err", err)
		r.ILogger.WithError(err).WithAny("dto", dto).Error(ctx, "error in send mediator create health check")

		return nil, err
	}

	return &dtos.HealthCheckDeleteResponse{
		Id: healthCheck.Id,
	}, nil
}
