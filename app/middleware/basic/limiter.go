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

func LimiterMiddleware(configs LimiterConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		max := 5
		if configs.Max > 0 {
			max = configs.Max
		}

		expiration := time.Minute
		if configs.Expiration > 0 {
			expiration = configs.Expiration
		}

		generator := fmt.Sprintf("%s:limit:%s:%s", config.Configs.Server.Name, ctx.Request.URL, ctx.ClientIP())
		if configs.keyGenerator != nil {
			generator = fmt.Sprintf("%s:limit:%s:%s", config.Configs.Server.Name, ctx.Request.URL, configs.keyGenerator(ctx))
		}

		fmt.Println(ctx.Request.URL.Path)

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
