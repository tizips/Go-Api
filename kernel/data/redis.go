package data

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"reflect"
	"saas/kernel/config"
	"time"
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

func StructToMap(items any) map[string]any {
	res := map[string]any{}
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

func Key(table string, id any) string {
	return fmt.Sprintf("%s:%s:%v", config.Values.Cache.Prefix, table, id)
}

func Ttl() time.Duration {
	return time.Duration(config.Values.Cache.Ttl) * time.Second
}
