package cache

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"reflect"
	"saas/app/helper"
	"saas/kernel/config"
	"saas/kernel/data"
	"strconv"
	"strings"
	"time"
)

type Model struct {
}

func (m *Model) AfterUpdate(tx *gorm.DB) (err error) {
	m.clear(tx)
	return
}

func (m *Model) AfterDelete(tx *gorm.DB) (err error) {
	m.clear(tx)
	return
}

//	数据修改之后，自动删除缓存模型
func (m *Model) clear(tx *gorm.DB) {

	var id interface{} = nil

	for _, field := range tx.Statement.Schema.Fields {
		if field.Name == tx.Statement.Schema.PrioritizedPrimaryField.Name {
			switch tx.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i < tx.Statement.ReflectValue.Len(); i++ {
					// 从字段中获取数值
					if fieldValue, isZero := field.ValueOf(tx.Statement.Context, tx.Statement.ReflectValue.Index(i)); !isZero {
						id = fieldValue
					}
				}
			case reflect.Struct:
				// 从字段中获取数值
				if fieldValue, isZero := field.ValueOf(tx.Statement.Context, tx.Statement.ReflectValue); !isZero {
					id = fieldValue
				}
			}
			break
		}
	}

	if id != nil {
		data.Redis.Del(tx.Statement.Context, key(tx.Statement.Schema.Table, id))
	}
}

//	优先从缓存中获取模型
func Find(ctx *gin.Context, id interface{}, model interface{}) {

	t := reflect.TypeOf(model).Elem()

	if t.Kind() == reflect.Struct {

		table := helper.StringSnake(t.Name())
		result, err := data.Redis.Get(ctx, key(table, id)).Result()

		if err == nil && result != "" {
			_ = json.Unmarshal([]byte(result), &model)
		}

		var ID interface{} = nil
		var keys = ""

		for i := 0; i < t.NumField(); i++ {
			tag := t.Field(i).Tag.Get("gorm")
			if tag != "" {
				tags := strings.Split(tag, ";")
				mark := false
				k := ""
				for _, item := range tags {
					if strings.HasPrefix(item, "column:") {
						k = strings.TrimPrefix(item, "column:")
					}
					if item == "primaryKey" || item == "primary_key" {
						mark = true
						break
					}
				}
				if mark {
					if k == "" {
						k = helper.StringSnake(t.Field(i).Name)
					}
					keys = k
					switch reflect.ValueOf(model).Elem().Field(i).Kind() {
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						n := reflect.ValueOf(model).Elem().Field(i).Int()
						if n > 0 {
							ID = n
						}
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						n := reflect.ValueOf(model).Elem().Field(i).Uint()
						if n > 0 {
							ID = n
						}
					case reflect.String:
						s := reflect.ValueOf(model).Elem().Field(i).String()
						if s != "" {
							ID = s
						}
					}
				}
			}
		}

		if ID == nil {
			tx := data.Database.Where(keys, id).First(&model)
			if tx.RowsAffected > 0 {
				hash, _ := json.Marshal(model)
				data.Redis.Set(ctx, key(table, id), string(hash), ttl())
			}
		}
	}
}

func key(table string, id interface{}) string {
	return fmt.Sprintf("%s:%s:%v", config.Configs.Cache.Prefix, table, id)
}

func ttl() time.Duration {
	t, err := strconv.Atoi(config.Configs.Cache.Ttl)
	if err != nil {
		t = 86400
	}
	return time.Duration(t) * time.Second
}
