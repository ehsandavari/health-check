package middlewares

import (
	"github.com/ehsandavari/go-jwt"
	"github.com/ehsandavari/go-logger"

	"health-check/infrastructure/config"
)

type Middleware struct {
	config     *config.SConfig
	logger     logger.ILogger
	iJwtServer jwt.IJwtServer
}

func NewMiddleware(config *config.SConfig, logger logger.ILogger, iJwtServer jwt.IJwtServer) *Middleware {
	return &Middleware{
		config:     config,
		logger:     logger,
		iJwtServer: iJwtServer,
	}
}
