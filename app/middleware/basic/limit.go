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

func LimitMiddleware(option *LimitOption) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		max := 5
		if option != nil && option.Max > 0 {
			max = option.Max
		}

		expiration := time.Minute
		if option != nil && option.Expiration > 0 {
			expiration = option.Expiration
		}

		generator := fmt.Sprintf("%s:limit:%s:%s", config.Values.Server.Name, ctx.Request.URL, ctx.ClientIP())
		if option != nil && option.keyGenerator != nil {
			generator = fmt.Sprintf("%s:limit:%s:%s", config.Values.Server.Name, ctx.Request.URL, option.keyGenerator(ctx))
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

type LimitOption struct {
	Max int

	Expiration time.Duration

	keyGenerator func(ctx *gin.Context) string
}
