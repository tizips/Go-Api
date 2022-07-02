package authorize

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"saas/app/constant"
)

func Jwt(ctx *gin.Context) *jwt.RegisteredClaims {

	var claims jwt.RegisteredClaims

	if data, exists := ctx.Get(constant.ContextJWT); exists {
		claims = data.(jwt.RegisteredClaims)
	} else {
		return nil
	}

	return &claims
}
