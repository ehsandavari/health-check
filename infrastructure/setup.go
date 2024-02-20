package infrastructure

import (
	"github.com/ehsandavari/go-context-plus"
	"github.com/ehsandavari/go-jwt"
	"github.com/ehsandavari/go-logger"
	"health-check/application/interfaces"
	"health-check/infrastructure/config"
	"health-check/infrastructure/cron"
	"health-check/infrastructure/notification"
	"health-check/infrastructure/postgres"
	"health-check/infrastructure/redis"
	"health-check/infrastructure/rest"
	"health-check/pkg/tracer"
	"time"
)

type Infrastructure struct {
	SConfig       *config.SConfig
	ILogger       logger.ILogger
	IJwtServer    jwt.IJwtServer
	ITracer       tracer.ITracer
	SPostgres     postgres.SPostgres
	IRedis        interfaces.IRedis
	ICron         interfaces.ICron
	IRest         interfaces.IRest
	INotification interfaces.INotification
}

func NewInfrastructure() *Infrastructure {
	sConfig := config.NewConfig()
	_logger := logger.NewLogger(
		*sConfig.Logger.IsDevelopment,
		*sConfig.Logger.DisableStacktrace,
		*sConfig.Logger.DisableStdout,
		sConfig.Logger.Level,
		sConfig.Service.Id,
		sConfig.Service.Name,
		sConfig.Service.Namespace,
		sConfig.Service.InstanceId,
		sConfig.Service.Version,
		sConfig.Service.Mode.String(),
		sConfig.Service.CommitId,
		logger.WithElk(sConfig.Logger.Elk.Url, sConfig.Logger.Elk.TimeoutSecond),
		logger.WithGormLogger(time.Duration(sConfig.Logger.Gorm.SlowThresholdMilliseconds)*time.Millisecond, *sConfig.Logger.Gorm.IgnoreRecordNotFoundError, *sConfig.Logger.Gorm.ParameterizedQueries),
	)
	_tracer := tracer.NewTracer(
		sConfig.Service.Name,
		sConfig.Tracer.Host,
		sConfig.Tracer.Port,
		_logger,
	)
	return &Infrastructure{
		SConfig: sConfig,
		ILogger: _logger,
		IJwtServer: jwt.NewJwtServer(
			sConfig.Jwt.Algorithm,
			sConfig.Jwt.PublicKey,
			sConfig.Jwt.PrivateKey,
			jwt.WithExpiresAt(time.Now().Add(time.Duration(sConfig.Jwt.ExpiresAtMinute)*time.Minute)),
			jwt.WithNotBefore(time.Now()),
			jwt.WithIssuedAt(time.Now()),
		),
		ITracer:       _tracer,
		SPostgres:     postgres.NewPostgres(sConfig.Postgres, _logger),
		IRedis:        redis.NewRedis(sConfig.Redis, _logger, _tracer),
		ICron:         cron.NewCron(_logger),
		IRest:         rest.NewRest(_logger),
		INotification: notification.NewNotification(sConfig.Notification, _logger, _tracer),
	}
}

func (r *Infrastructure) Close() {
	ctx := contextplus.Background()

	if err := r.SPostgres.Close(); err != nil {
		r.ILogger.WithError(err).Error(ctx, "error in close postgres")
	}

	if err := r.IRedis.Close(); err != nil {
		r.ILogger.WithError(err).Error(ctx, "error in close redis")
	}

	if err := r.ITracer.Close(); err != nil {
		r.ILogger.WithError(err).Error(ctx, "error in tracer close")
	}

	if err := r.ILogger.Sync(); err != nil {
		r.ILogger.WithError(err).Error(ctx, "error in sync logger")
	}
}
