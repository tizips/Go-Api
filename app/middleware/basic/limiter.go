package basic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"saas/kernel/config"
	"saas/kernel/data"
	"saas/kernel/response"
	"time"
)

func LimiterMiddleware(configs ...LimiterConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		max := 5
		if len(configs) > 0 && configs[0].Max > 0 {
			max = configs[0].Max
		}

		expiration := time.Minute
		if len(configs) > 0 && configs[0].Expiration > 0 {
			expiration = configs[0].Expiration
		}

		generator := fmt.Sprintf("%s:limit:%s:%s", config.Values.Server.Name, ctx.Request.URL, ctx.ClientIP())
		if len(configs) > 0 && configs[0].keyGenerator != nil {
			generator = fmt.Sprintf("%s:limit:%s:%s", config.Values.Server.Name, ctx.Request.URL, configs[0].keyGenerator(ctx))
		}

		number, err := data.Redis.Incr(ctx, generator).Result()
		if err != nil || number > int64(max) {
			ctx.Abort()
			ctx.JSON(http.StatusOK, response.Response{
				Code:    50000,
				Message: "访问受限，请稍后重试",
			})
			return
		}

		data.Redis.Expire(ctx, generator, expiration)

		ctx.Next()
	}
}

type LimiterConfig struct {
	Max int

	Expiration time.Duration

	keyGenerator func(ctx *gin.Context) string
}
