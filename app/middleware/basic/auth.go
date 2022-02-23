package basic

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"saas/kernel/auth"
	"saas/kernel/response"
)

func Auth(issue string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if !auth.Check(ctx) {
			ctx.Abort()
			ctx.JSON(http.StatusUnauthorized, response.Response{
				Code:    40100,
				Message: "Unauthorized1",
			})
			return
		}

		claims := ctx.MustGet("JWT").(jwt.StandardClaims)

		if !claims.VerifyIssuer(issue, true) {
			ctx.Abort()
			ctx.JSON(http.StatusUnauthorized, response.Response{
				Code:    40100,
				Message: "Unauthorized2",
			})
			return
		}

		ctx.Next()
	}
}
