package middlewares

import (
	"github.com/ehsandavari/go-context-plus"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"health-check/pkg/apiHandler"
	"net/http"
	"strings"
)

const _httpStatusUnauthorized = http.StatusUnauthorized

var _httpStatusUnauthorizedText = http.StatusText(_httpStatusUnauthorized)

func (r *Middleware) Jwt() gin.HandlerFunc {
	return func(ctxGin *gin.Context) {
		reqCtx := ctxGin.Request.Context()
		ctx := contextplus.FromContext(reqCtx)

		authorization := ctxGin.GetHeader("authorization")
		if len(authorization) == 0 {
			r.logger.Warn(ctx, "authorization not set in request header")
			ctxGin.AbortWithStatusJSON(_httpStatusUnauthorized, apiHandler.NewBaseApiResponse(false, apiHandler.NewApiError(_httpStatusUnauthorized, _httpStatusUnauthorizedText)))
			return
		}

		token := strings.TrimPrefix(authorization, "Bearer ")
		if len(token) == 0 || token == authorization {
			r.logger.Warn(ctx, "authorization token is invalid value")
			ctxGin.AbortWithStatusJSON(_httpStatusUnauthorized, apiHandler.NewBaseApiResponse(false, apiHandler.NewApiError(_httpStatusUnauthorized, _httpStatusUnauthorizedText)))
			return
		}

		valid, err := r.iJwtServer.VerifyToken(token, "", "")
		if err != nil {
			r.logger.WithError(err).Error(ctx, "error in Verify authorization token")
			ctxGin.AbortWithStatusJSON(_httpStatusUnauthorized, apiHandler.NewBaseApiResponse(false, apiHandler.NewApiError(_httpStatusUnauthorized, _httpStatusUnauthorizedText)))
			return
		}

		if !valid {
			r.logger.Warn(ctx, "authorization token is invalid")
			ctxGin.AbortWithStatusJSON(_httpStatusUnauthorized, apiHandler.NewBaseApiResponse(false, apiHandler.NewApiError(_httpStatusUnauthorized, _httpStatusUnauthorizedText)))
			return
		}

		ctx.User.SetId(uuid.MustParse(r.iJwtServer.GetUserId()))
		if len(r.iJwtServer.GetEmail()) != 0 {
			ctx.User.SetEmail(r.iJwtServer.GetEmail())
			ctx.User.SetEmailVerified(r.iJwtServer.GetEmailVerified())
		}
		if len(r.iJwtServer.GetPhoneNumber()) != 0 {
			ctx.User.SetPhoneNumber(r.iJwtServer.GetPhoneNumber())
			ctx.User.SetPhoneNumberVerified(r.iJwtServer.GetPhoneNumberVerified())
		}

		ctxGin.Request = ctxGin.Request.WithContext(ctx.ToContext())
		ctxGin.Next()
	}
}
