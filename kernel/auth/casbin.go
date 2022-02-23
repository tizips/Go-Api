package auth

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	adapter "github.com/casbin/gorm-adapter/v3"
	"saas/kernel/config"
	"saas/kernel/data"
)

const ROOT = 888

var Casbin *casbin.Enforcer

func InitCasbin() {

	a, err := adapter.NewAdapterByDBUseTableName(data.Database, config.Configs.Database.Prefix, "sys_casbin")
	if err != nil {
		fmt.Printf("Casbin new adapter error:%v", err)
		return
	}

	Casbin, err = casbin.NewEnforcer("conf/casbin.conf", a)
	if err != nil {
		fmt.Printf("Casbin new enforcer error:%v", err)
		return
	}

}

func NameByAdmin(id interface{}) string {
	return fmt.Sprintf("admin:%v", id)
}

func NameByRole(id interface{}) string {
	return fmt.Sprintf("role:%v", id)
}

func Root(id interface{}) bool {
	exist, _ := Casbin.HasRoleForUser(NameByAdmin(id), NameByRole(ROOT))
	return exist
}
