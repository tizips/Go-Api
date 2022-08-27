package basic

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"saas/kernel/app"
	"saas/kernel/authorize"
	"time"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		start := time.Now()

		ctx.Next()

		id := authorize.Id(ctx)

		fields := logrus.Fields{
			"status":  ctx.Writer.Status(),
			"method":  ctx.Request.Method,
			"uri":     ctx.FullPath(),
			"path":    ctx.Request.URL.Path,
			"ip":      ctx.ClientIP(),
			"latency": time.Now().Sub(start),
		}

		if id > 0 {
			fields["admin"] = id
		}

		app.Logger.Api.Info(fields)

	}
}
