package basic

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-module/carbon/v2"
	"saas/kernel/config"
	"saas/kernel/data"
	"time"
)

func CheckJwt(ctx context.Context, channel string, claims jwt.RegisteredClaims) bool {

	result, _ := data.Redis.Exists(ctx, Blacklist(channel, "login", claims.Subject)).Result()

	return result != 1
}

func BlackJwt(ctx context.Context, channel string, claims jwt.RegisteredClaims) bool {

	now := carbon.Now()

	expires := time.Duration(claims.ExpiresAt.Unix()+12*60*60-now.Timestamp()) * time.Second

	_, err := data.Redis.Set(ctx, Blacklist(channel, "login", claims.Subject), now.ToDateTimeString(), expires).Result()

	if err != nil {
		return false
	}

	return true
}

func Blacklist(channel string, method string, str string) string {
	return config.Values.Server.Name + ":" + channel + ":blacklist:" + method + ":" + str
}
