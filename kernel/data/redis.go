package data

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"reflect"
	"saas/kernel/config"
	"strconv"
	"time"
)

var Redis *redis.Client

func InitRedis() {

	db, err := strconv.Atoi(config.Configs.Redis.Db)
	if err != nil {
		db = 0
	}

	Redis = redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%s", config.Configs.Redis.Host, config.Configs.Redis.Port),
		Password:    config.Configs.Redis.Auth,
		DB:          db,
		MaxConnAge:  100,
		PoolTimeout: 3,
		IdleTimeout: 60,
	})
}

func StructToMap(items interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	if items == nil {
		return res
	}
	v := reflect.TypeOf(items)
	reflectValue := reflect.Indirect(reflect.ValueOf(items))
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Name
		if tag != "" && v.Field(i).Type.Kind() != reflect.Struct && v.Field(i).Type.Kind() != reflect.Slice {
			field := reflectValue.Field(i).Interface()
			res[tag] = field
		}
	}
	return res
}

func Key(table string, id interface{}) string {
	return fmt.Sprintf("%s:%s:%v", config.Configs.Cache.Prefix, table, id)
}

func Ttl() time.Duration {
	ttl, err := strconv.Atoi(config.Configs.Cache.Ttl)
	if err != nil {
		ttl = 86400
	}
	return time.Duration(ttl) * time.Second
}
