package notification

import (
	"github.com/ehsandavari/go-context-plus"
	"github.com/ehsandavari/go-logger"
	"github.com/nikoksr/notify"
	"health-check/application/interfaces"
	"health-check/pkg/tracer"
)

type sNotification struct {
	logger logger.ILogger
	tracer tracer.ITracer
	config *SConfig
	notify *notify.Notify
}

func NewNotification(config *SConfig, logger logger.ILogger, tracer tracer.ITracer) interfaces.INotification {
	n := sNotification{
		logger: logger,
		tracer: tracer,
		config: config,
		notify: notify.New(),
	}
	n.AddDiscord()
	n.AddSlack()
	return n
}

func (r sNotification) Send(ctx *contextplus.Context, subject string, message string) error {
	span, ctx := r.tracer.SpanFromContext(ctx)
	defer span.Finish()

	if err := r.notify.Send(ctx.Context, subject, message); err != nil {
		span.SetTag("error", true)
		span.LogKV("err", err)
		r.logger.WithError(err).WithString("subject", subject).WithString("message", message).Error(ctx, "error in send notification")

		return err
	}

	return nil
}
