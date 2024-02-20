package middlewares

import (
	"github.com/ehsandavari/go-context-plus"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func (r *Middleware) RequestId() gin.HandlerFunc {
	return requestid.New(
		requestid.WithHandler(
			func(ctxGin *gin.Context, requestID string) {
				reqCtx := ctxGin.Request.Context()
				ctx := contextplus.FromContext(reqCtx)
				ctx.SetRequestId(requestID)
				ctxGin.Request = ctxGin.Request.WithContext(ctx.ToContext())
				ctxGin.Next()
			},
		),
	)
}
