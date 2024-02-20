package grpc

import (
	"github.com/ehsandavari/go-context-plus"
	"github.com/ehsandavari/go-logger"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"health-check/application"
	"health-check/infrastructure/config"
	"health-check/pkg/tracer"
	//"healthCheck/presentation/grpc/proto/aaa"
	//"healthCheck/presentation/grpc/services"
	"log"
	"net"
)

type SGrpc struct {
	server  *grpc.Server
	sConfig *config.SConfig
	iLogger logger.ILogger
	iTracer tracer.ITracer
}

func NewSGrpc(application *application.Application, sConfig *config.SConfig, iLogger logger.ILogger, iTracer tracer.ITracer) *SGrpc {
	var sGrpc SGrpc
	sGrpc.sConfig = sConfig
	if *sConfig.Service.Grpc.IsEnabled {
		sGrpc.server = grpc.NewServer()
		sGrpc.iLogger = iLogger
		sGrpc.iTracer = iTracer
	}
	return &sGrpc
}

func (r *SGrpc) Start() {
	if *r.sConfig.Service.Grpc.IsEnabled {
		netListener, err := net.Listen("tcp", ":"+r.sConfig.Service.Grpc.Port)
		if err != nil {
			log.Fatal("error in net listen ", err)
		}
		//aaa.RegisterAuthServiceServer(r.server, services.NewAuthService())
		grpcPrometheus.Register(r.server)

		if *r.sConfig.Service.Grpc.IsDevelopment {
			reflection.Register(r.server)
		}

		ctx := contextplus.Background()

		go func() {
			if err = r.server.Serve(netListener); err != nil {
				r.iLogger.WithError(err).Fatal(ctx, "error in serve grpc server")
			}
		}()
		r.iLogger.WithAny("grpcServerInfo", r.sConfig.Service.Grpc).Info(ctx, "grpc server start")
	}
}

func (r *SGrpc) Stop() {
	if *r.sConfig.Service.Grpc.IsEnabled {
		r.server.GracefulStop()
	}
}
