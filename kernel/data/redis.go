package data

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"saas/kernel/config"
)

var Redis *redis.Client

func InitRedis() {

	Redis = redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%d", config.Values.Redis.Host, config.Values.Redis.Port),
		Password:    config.Values.Redis.Auth,
		DB:          config.Values.Redis.Db,
		MaxConnAge:  100,
		PoolTimeout: 3,
		IdleTimeout: 60,
	})
}
