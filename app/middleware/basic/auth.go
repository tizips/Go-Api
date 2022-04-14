package basic

import (
	"github.com/gin-gonic/gin"
	"saas/app/service/basic"
	"saas/kernel/auth"
	"saas/kernel/response"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if !auth.Check(ctx) {
			ctx.Abort()
			response.ToResponseByUnauthorized(ctx)
			return
		}

		claims := auth.Jwt(ctx)

		if !claims.VerifyIssuer("admin", true) {
			ctx.Abort()
			response.ToResponseByUnauthorized(ctx)
			return
		}

		if !basic.CheckJwt(ctx, "admin", *claims) {
			ctx.Abort()
			response.ToResponseByUnauthorized(ctx)
			return
		}

		ctx.Next()
	}
}
