package basic

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"saas/app/service/helper"
	"saas/kernel/config"
	"saas/kernel/data"
	"strconv"
)

func JwtParseMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if authorization := ctx.GetHeader("Authorization"); authorization != "" {
			var claims jwt.StandardClaims
			_, _ = jwt.ParseWithClaims(authorization, &claims, func(token *jwt.Token) (any, error) {
				return []byte(config.Configs.Jwt.Secret), nil
			})
			if claims.Id != "" {
				now := carbon.Now()
				if ok := claims.VerifyNotBefore(now.Timestamp(), false); ok {
					if ok = claims.VerifyExpiresAt(now.Timestamp(), true); ok { //	生效的授权令牌操作
						set(ctx, claims)
					} else if ok = claims.VerifyExpiresAt(now.SubHours(12).Timestamp(), true); ok { //	失效的授权令牌，重新发放
						refresh(ctx, claims)
					}
				}

				ctx.Next()
			}
		}
	}
}

func set(ctx *gin.Context, claims jwt.StandardClaims) {
	ctx.Set("ID", claims.Id)
	ctx.Set("JWT", claims)
}

func refresh(ctx *gin.Context, claims jwt.StandardClaims) {

	cache, _ := data.Redis.HGetAll(ctx, refreshKey(claims.Audience)).Result()

	now := carbon.Now()

	if len(cache) <= 0 {

		expires := claims.ExpiresAt
		audience := claims.Audience

		claims.NotBefore = now.Timestamp()
		claims.IssuedAt = now.Timestamp()
		claims.ExpiresAt = now.AddHours(12).Timestamp()
		claims.Audience = helper.JwtToken(claims.Id)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		if signed, err := token.SignedString([]byte(config.Configs.Jwt.Secret)); err == nil {

			set(ctx, claims)

			ctx.Header("Authorization", signed)

			affected, err := data.Redis.HSet(ctx, refreshKey(audience), "token", signed, "created_at", now.ToDateTimeString()).Result()
			if err == nil && affected > 0 {
				data.Redis.ExpireAt(ctx, refreshKey(audience), carbon.CreateFromTimestamp(expires).AddHours(12).Carbon2Time())
			}
		}

	} else {

		leeway, _ := strconv.Atoi(config.Configs.Jwt.Leeway)

		diff := now.DiffInSecondsWithAbs(carbon.Parse(cache["created_at"]))

		if diff <= int64(leeway) {
			ctx.Header("Authorization", cache["token"])
		}
	}
}

func refreshKey(token string) string {
	return fmt.Sprintf("saas:token:refresh:%s", token)
}
