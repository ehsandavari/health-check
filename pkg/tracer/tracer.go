package tracer

import (
	"github.com/ehsandavari/go-context-plus"
	"github.com/ehsandavari/go-logger"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"io"
	"runtime"
	"strings"
)

//go:generate mockgen -destination=./tracer_mock.go -package=tracer . ITracer,ISpan

type ITracer interface {
	GetTracer() opentracing.Tracer
	SpanFromContext(ctx *contextplus.Context, opts ...opentracing.StartSpanOption) (ISpan, *contextplus.Context)
	Close() error
}

type ISpan interface {
	opentracing.Span
}

type sTracer struct {
	serviceName string
	tracer      opentracing.Tracer
	closer      io.Closer
}

func NewTracer(serviceName string, host string, Port string, logger logger.ILogger) ITracer {
	cfg := jaegerConfig.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LocalAgentHostPort: host + ":" + Port,
		},
	}

	t, closer, err := cfg.NewTracer(jaegerConfig.ZipkinSharedRPCSpan(true))
	if err != nil {
		logger.WithError(err).Fatal(contextplus.Background(), "error in during Listen jaeger")
	}

	opentracing.SetGlobalTracer(t)

	return &sTracer{
		serviceName: serviceName,
		tracer:      t,
		closer:      closer,
	}
}

func (r *sTracer) GetTracer() opentracing.Tracer {
	return r.tracer
}

func (r *sTracer) SpanFromContext(ctx *contextplus.Context, opts ...opentracing.StartSpanOption) (span ISpan, _ *contextplus.Context) {
	span, ctx.Context = opentracing.StartSpanFromContext(ctx.Context, r.getCallerInfo(), opts...)
	if ctx.User.Id() != uuid.Nil {
		span.SetTag("user_id", ctx.User.Id())
	}
	if len(ctx.User.PhoneNumber()) != 0 {
		span.SetTag("phone_number", ctx.User.PhoneNumber())
	}
	span.SetTag("request_id", ctx.RequestId())
	if sc, ok := span.Context().(jaeger.SpanContext); ok {
		ctx.SetTraceId(sc.TraceID().String())
	}
	return span, ctx
}

func (r *sTracer) getCallerInfo() string {
	pc, _, _, _ := runtime.Caller(2)
	callerFunc := runtime.FuncForPC(pc)
	caller := "unknown"

	if callerFunc != nil {
		caller = strings.ReplaceAll(callerFunc.Name(), r.serviceName+"/", "")
	}

	return caller
}

func (r *sTracer) Close() error {
	return r.closer.Close()
}
