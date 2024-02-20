package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ehsandavari/go-context-plus"
	"github.com/ehsandavari/go-logger"
	"github.com/redis/go-redis/v9"
	"health-check/application/interfaces"
	"health-check/pkg/tracer"
	"time"
)

type sRedis struct {
	serviceName string
	logger      logger.ILogger
	tracer      tracer.ITracer
	client      *redis.Client
}

func NewRedis(config *SConfig, logger logger.ILogger, tracer tracer.ITracer) interfaces.IRedis {
	return &sRedis{
		logger: logger,
		tracer: tracer,
		client: redis.NewClient(&redis.Options{
			Addr: config.Host + ":" + config.Port,
			DB:   0,
		}),
	}
}

func (r *sRedis) Set(ctx *contextplus.Context, key string, value any) error {
	span, ctx := r.tracer.SpanFromContext(ctx)
	defer span.Finish()

	valueByte, err := json.Marshal(value)
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("err", err)
		r.logger.WithError(err).WithAny("value", value).Error(ctx, "error in json marshal")

		return err
	}

	if err := r.client.Set(ctx, fmt.Sprintf("%s:%s", r.serviceName, key), valueByte, time.Minute*5).Err(); err != nil {
		span.SetTag("error", true)
		span.LogKV("err", err)
		r.logger.WithError(err).WithString("key", fmt.Sprintf("%s:%s", r.serviceName, key)).WithAny("value", value).WithByteString("valueByte", valueByte).Error(ctx, "error in redis set")

		return err
	}

	return nil
}

func (r *sRedis) Get(ctx *contextplus.Context, key string) (string, error) {
	span, ctx := r.tracer.SpanFromContext(ctx)
	defer span.Finish()

	val, err := r.client.Get(ctx, fmt.Sprintf("%s:%s", r.serviceName, key)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		span.SetTag("error", true)
		span.LogKV("err", err)
		r.logger.WithError(err).Error(ctx, "error in redis get")

		return "", err
	}

	return val, nil
}

func (r *sRedis) Publish(ctx *contextplus.Context, channelName string, message any) error {
	span, ctx := r.tracer.SpanFromContext(ctx)
	defer span.Finish()

	if err := r.client.Publish(ctx, channelName, message).Err(); err != nil {
		span.SetTag("error", true)
		span.LogKV("err", err)
		r.logger.WithError(err).Error(ctx, "error in redis publish")
		return err
	}

	return nil
}

func (r *sRedis) Subscribe(ctx *contextplus.Context, channelName string, channel chan<- string) {
	span, ctx := r.tracer.SpanFromContext(ctx)
	defer span.Finish()

	var subscribe = r.client.Subscribe(ctx, channelName)
	defer subscribe.Close()

	for msg := range subscribe.Channel() {
		channel <- msg.Payload
	}
}

func (r *sRedis) Close() error {
	return r.client.Close()
}
