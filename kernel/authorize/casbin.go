package authorize

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gookit/color"
	"saas/kernel/app"
)

const ROOT = 666

func InitCasbin() {

	a, err := adapter.NewAdapterByDBUseTableName(app.Database, app.Cfg.Database.MySQL.Prefix, "sys_casbin")
	if err != nil {
		color.Errorf("Casbin new adapter error: %v", err)
		return
	}

	app.Casbin, err = casbin.NewEnforcer(app.Dir.Root+"/conf/casbin.conf", a)
	if err != nil {
		color.Errorf("Casbin new enforcer error: %v", err)
		return
	}

}

func NameByAdmin(id any) string {
	return fmt.Sprintf("admin:%v", id)
}

func NameByRole(id any) string {
	return fmt.Sprintf("role:%v", id)
}

func Root(id any) bool {
	exist, _ := app.Casbin.HasRoleForUser(NameByAdmin(id), NameByRole(ROOT))
	return exist
}
