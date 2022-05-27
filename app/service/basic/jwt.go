package basic

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang-module/carbon/v2"
	"saas/kernel/config"
	"saas/kernel/data"
	"time"
)

func CheckJwt(ctx context.Context, channel string, claims jwt.StandardClaims) bool {

	result, _ := data.Redis.Exists(ctx, Blacklist(channel, "login", claims.Audience)).Result()

	return result != 1
}

func BlackJwt(ctx context.Context, channel string, claims jwt.StandardClaims) bool {

	now := carbon.Now()

	expires := time.Duration(claims.ExpiresAt+12*60*60-now.Timestamp()) * time.Second

	_, err := data.Redis.Set(ctx, Blacklist(channel, "login", claims.Audience), now.ToDateTimeString(), expires).Result()

	if err != nil {
		return false
	}

	return true
}

func Blacklist(channel string, method string, str string) string {
	return config.Values.Redis.Prefix + ":" + channel + ":blacklist:" + method + ":" + str
}
