package database

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"saas/kernel/app"
	"time"
)

func InitRedis() {

	app.Redis = redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%d", app.Cfg.Database.Redis.Host, app.Cfg.Database.Redis.Port),
		Password:    app.Cfg.Database.Redis.Password,
		DB:          app.Cfg.Database.Redis.Db,
		MaxConnAge:  time.Duration(app.Cfg.Database.Redis.MaxConnAge),
		PoolTimeout: time.Duration(app.Cfg.Database.Redis.PoolTimeout),
		IdleTimeout: time.Duration(app.Cfg.Database.Redis.IdleTimeout),
	})
}
