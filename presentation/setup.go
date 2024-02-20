package presentation

import (
	"health-check/application"
	"health-check/infrastructure"
	"health-check/presentation/api"
	"health-check/presentation/grpc"
)

type Presentation struct {
	sApi  *api.SApi
	sGrpc *grpc.SGrpc
}

func NewPresentation(infrastructure *infrastructure.Infrastructure, application *application.Application) *Presentation {
	return &Presentation{
		sApi:  api.NewSApi(application, infrastructure.SConfig, infrastructure.IJwtServer, infrastructure.ILogger, infrastructure.ITracer),
		sGrpc: grpc.NewSGrpc(application, infrastructure.SConfig, infrastructure.ILogger, infrastructure.ITracer),
	}
}

func (r *Presentation) Setup() {
	r.sApi.Start()
	r.sGrpc.Start()
}

func (r *Presentation) Close() {
	r.sApi.Stop()
	r.sGrpc.Stop()
}
