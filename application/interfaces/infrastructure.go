package interfaces

import (
	"github.com/ehsandavari/go-context-plus"
	"health-check/domain/enums"
	"net/http"
	"time"
)

//go:generate mockgen -destination=./infrastructure_mock.go -package=interfaces . ICron,INotification,IRedis,IRest

type ICron interface {
	AddJob(key uint, createAt time.Time, interval string, job func()) error
	RemoveJob(key uint)
}

type INotification interface {
	Send(ctx *contextplus.Context, subject string, message string) error
}

type IRedis interface {
	Set(ctx *contextplus.Context, key string, value any) error
	Get(ctx *contextplus.Context, key string) (string, error)
	Publish(ctx *contextplus.Context, channelName string, message any) error
	Subscribe(ctx *contextplus.Context, channelName string, channel chan<- string)
	Close() error
}

type IRest interface {
	Execute(ctx *contextplus.Context, method enums.HttpMethod, url string, headers map[string]string, body map[string]any) (int, http.Header, string, error)
}
