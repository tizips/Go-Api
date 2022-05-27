package authorize

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"saas/app/constant"
)

func Jwt(ctx *gin.Context) *jwt.StandardClaims {

	var claims jwt.StandardClaims

	if data, exists := ctx.Get(constant.ContextJWT); exists {
		claims = data.(jwt.StandardClaims)
	} else {
		return nil
	}

	return &claims
}
