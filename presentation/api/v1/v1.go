package v1

import (
	"github.com/ehsandavari/go-logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"health-check/application"
	"health-check/pkg/tracer"
	"health-check/presentation/api/middlewares"
	"health-check/presentation/api/v1/controllers"
)

//	@title			api
//	@version		1.0
//	@description	Example Api

//	@contact.name	Ehsan Davari
//	@contact.url	https://github.com/ehsandavari
//	@contact.email	ehsandavari.ir@gmail.com

//	@BasePath	/api/v1

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImVoc2FuZGF2YXJpLmlyQGdtYWlsLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwicGhvbmVfbnVtYmVyIjoiMDkyMTU1ODA2OTAiLCJwaG9uZV9udW1iZXJfdmVyaWZpZWQiOnRydWUsImlzcyI6IldpdGhJc3N1ZXIiLCJzdWIiOiJlNGRhMDhlMS05NmQ0LTQ1NTgtOWZiOS1jN2UwNGJiMzdlMDIiLCJhdWQiOlsiYXBpMSIsImFwaTIiXSwiZXhwIjoxNzI1NjY1ODE0LCJuYmYiOjE2OTQxMDgyMTQsImlhdCI6MTY5NDEwODIxNCwianRpIjoiYXNkbG1rc2ZrZGZtYWtzZGZtYXNrbGQifQ.GECrzoQ2vEDvS-wGRII62BjtNmxDGDo38stQb2_-t1IGfIIP0hr3C_iHWN7odXx3r6HrkMuyyO2gFEdkd17qRSfFsWi-4oJnzVvnzRYA1uIAnmg9QplrHNpb4mbedC0BpZSkcoju-3hNFoeDuc1_0ZJyEMHHexhb64Jou1XzVss

type V1 struct {
	application *application.Application
	routerGroup *gin.RouterGroup
	middleware  *middlewares.Middleware
	iLogger     logger.ILogger
	iTracer     tracer.ITracer
}

func NewV1(application *application.Application, routerGroup *gin.RouterGroup, middleware *middlewares.Middleware, iLogger logger.ILogger, iTracer tracer.ITracer) *V1 {
	return &V1{
		application: application,
		routerGroup: routerGroup,
		middleware:  middleware,
		iLogger:     iLogger,
		iTracer:     iTracer,
	}
}

func (r *V1) Setup() {
	apiRouterGroup := r.routerGroup.Group("/v1")
	{
		apiRouterGroup.GET("/swagger/*any", ginSwagger.WrapHandler(
			swaggerFiles.NewHandler(),
			ginSwagger.InstanceName("v1"),
		))

		controllers.NewHealthCheckController(r.application, apiRouterGroup, r.iLogger, r.iTracer)

		apiRouterGroup.Use(r.middleware.Jwt())
		{
			//controllers.NewUserController(apiRouterGroup, r.iLogger, r.iTracer)
		}
	}
}
