package cache

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"reflect"
	"saas/app/helper/str"
	"saas/kernel/app"
	"strings"
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

// 数据修改之后，自动删除缓存模型
func (m *Model) clear(tx *gorm.DB) {

	key := id(tx)

	if key != "" {
		app.Redis.Del(tx.Statement.Context, Key(tx.Statement.Schema.Table, key))
	}
}

// 优先从缓存中获取模型
func FindById(ctx *gin.Context, model any, id any) {

	t := reflect.TypeOf(model).Elem()

	if t.Kind() == reflect.Struct {

		table := str.Snake(t.Name())
		result, err := app.Redis.Get(ctx, Key(table, id)).Result()

		if err == nil && result != "" {
			_ = json.Unmarshal([]byte(result), &model)
			return
		}

		tx := app.MySQL.Find(&model, id)
		if tx.RowsAffected > 0 {
			hash, _ := json.Marshal(model)
			app.Redis.Set(ctx, Key(table, id), string(hash), ttl())
		}
	}
}

func id(tx *gorm.DB) string {

	var ids = make([]string, 0)

	for _, field := range tx.Statement.Schema.Fields {
		if field.Name == tx.Statement.Schema.PrioritizedPrimaryField.Name {
			switch tx.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i < tx.Statement.ReflectValue.Len(); i++ {
					// 从字段中获取数值
					if fieldValue, isZero := field.ValueOf(tx.Statement.Context, tx.Statement.ReflectValue.Index(i)); !isZero {
						ids = append(ids, fmt.Sprintf("%v", fieldValue))
					}
				}
			case reflect.Struct:
				// 从字段中获取数值
				if fieldValue, isZero := field.ValueOf(tx.Statement.Context, tx.Statement.ReflectValue); !isZero {
					ids = append(ids, fmt.Sprintf("%v", fieldValue))
				}
			}
		}
	}

	return strings.Join(ids, "-")
}
