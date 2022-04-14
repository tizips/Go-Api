package basic

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"saas/app/constant"
	"saas/app/service/helper"
	"saas/kernel/config"
	"saas/kernel/data"
)

func JwtParseMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if authorization := ctx.GetHeader("Authorization"); authorization != "" {
			var claims jwt.StandardClaims
			_, _ = jwt.ParseWithClaims(authorization, &claims, func(token *jwt.Token) (any, error) {
				return []byte(config.Values.Jwt.Secret), nil
			})
			if claims.Id != "" {
				now := carbon.Now()
				if ok := claims.VerifyNotBefore(now.Timestamp(), false); ok {
					if ok = claims.VerifyExpiresAt(now.Timestamp(), true); ok { //	生效的授权令牌操作
						set(ctx, claims)
					} else if ok = claims.VerifyExpiresAt(now.SubHours(config.Values.Jwt.Lifetime).Timestamp(), true); ok { //	失效的授权令牌，重新发放
						refresh(ctx, claims)
					}
				}

				ctx.Next()
			}
		}
	}
}

func set(ctx *gin.Context, claims jwt.StandardClaims) {
	ctx.Set(constant.ContextID, claims.Id)
	ctx.Set(constant.ContextJWT, claims)
}

func refresh(ctx *gin.Context, claims jwt.StandardClaims) {

	cache, _ := data.Redis.HGetAll(ctx, refreshKey(claims.Audience)).Result()

	now := carbon.Now()

	if len(cache) <= 0 {

		expires := claims.ExpiresAt
		audience := claims.Audience

		claims.NotBefore = now.Timestamp()
		claims.IssuedAt = now.Timestamp()
		claims.ExpiresAt = now.AddHours(config.Values.Jwt.Lifetime).Timestamp()
		claims.Audience = helper.JwtToken(claims.Id)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		if signed, err := token.SignedString([]byte(config.Values.Jwt.Secret)); err == nil {

			set(ctx, claims)

			ctx.Header("Authorization", signed)

			affected, err := data.Redis.HSet(ctx, refreshKey(audience), "token", signed, "created_at", now.ToDateTimeString()).Result()
			if err == nil && affected > 0 {
				data.Redis.ExpireAt(ctx, refreshKey(audience), carbon.CreateFromTimestamp(expires).AddHours(config.Values.Jwt.Lifetime).Carbon2Time())
			}
		}

	} else {

		diff := now.DiffInSecondsWithAbs(carbon.Parse(cache["created_at"]))

		if diff <= config.Values.Jwt.Leeway {
			ctx.Header("Authorization", cache["token"])
		}
	}
}

func refreshKey(token string) string {
	return fmt.Sprintf("saas:token:refresh:%s", token)
}
