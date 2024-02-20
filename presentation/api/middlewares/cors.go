package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (r *Middleware) Cors() gin.HandlerFunc {
	return cors.Default()
}
