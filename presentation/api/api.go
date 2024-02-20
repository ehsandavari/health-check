package api

import (
	"github.com/ehsandavari/go-context-plus"
	"github.com/ehsandavari/go-jwt"
	"github.com/ehsandavari/go-logger"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"health-check/application"
	"health-check/infrastructure/config"
	"health-check/pkg/tracer"
	_ "health-check/presentation/api/docs"
	"health-check/presentation/api/middlewares"
	"health-check/presentation/api/v1"
	"net/http"
)

type SApi struct {
	application *application.Application
	server      *http.Server
	sConfig     *config.SConfig
	iJwtServer  jwt.IJwtServer
	iLogger     logger.ILogger
	iTracer     tracer.ITracer
}

func NewSApi(application *application.Application, sConfig *config.SConfig, iJwtServer jwt.IJwtServer, iLogger logger.ILogger, iTracer tracer.ITracer) *SApi {
	var sApi SApi
	sApi.sConfig = sConfig
	if *sConfig.Service.Api.IsEnabled {
		sApi.application = application
		sApi.server = &http.Server{
			Addr: sConfig.Service.Api.Host + ":" + sConfig.Service.Api.Port,
		}
		sApi.iJwtServer = iJwtServer
		sApi.iLogger = iLogger
		sApi.iTracer = iTracer
	}
	return &sApi
}

type iValidator interface {
	IsValid() bool
}

func Validator(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(iValidator)
	return value.IsValid()
}

func (r *SApi) Start() {
	if *r.sConfig.Service.Api.IsEnabled {
		gin.SetMode(r.sConfig.Service.Api.Mode)

		engine := gin.Default()

		ctx := contextplus.Background()

		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			if err := v.RegisterValidation("enum", Validator); err != nil {
				r.iLogger.WithError(err).Fatal(ctx, "error in register validation")
			}
		}

		middleware := middlewares.NewMiddleware(r.sConfig, r.iLogger, r.iJwtServer)
		engine.Use(
			middleware.Cors(),
			middleware.I18n(),
			middleware.RequestId(),
		)

		monitoringRouterGroup := engine.Group("/-")
		{
			monitoringRouterGroup.GET("/health", func(ctx *gin.Context) { ctx.Status(http.StatusOK) })
			monitoringRouterGroup.GET("/liveness", func(ctx *gin.Context) { ctx.Status(http.StatusOK) })
			monitoringRouterGroup.GET("/readiness", func(ctx *gin.Context) { ctx.Status(http.StatusOK) })
			monitoringRouterGroup.GET("/metrics", gin.WrapH(promhttp.Handler()))
		}

		apiRouterGroup := engine.Group("/api")
		{
			v1.NewV1(r.application, apiRouterGroup, middleware, r.iLogger, r.iTracer).Setup()
		}

		go func() {
			r.server.Handler = engine.Handler()
			if err := r.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				r.iLogger.WithError(err).Fatal(ctx, "error in serve api server")
			}
		}()
		r.iLogger.WithAny("api server info", r.sConfig.Service.Api).Info(ctx, "api server start")
	}
}

func (r *SApi) Stop() {
	if *r.sConfig.Service.Api.IsEnabled {
		ctx := contextplus.Background()
		if err := r.server.Shutdown(ctx); err != nil {
			r.iLogger.WithError(err).Error(ctx, "error in shutdown api server")
		}
	}
}
