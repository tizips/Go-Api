package helper

import (
	"fmt"
	"github.com/golang-module/carbon/v2"
	"saas/app/helper"
	"saas/kernel/config"
)

func JwtToken(id interface{}) string {
	now := carbon.Now()
	return helper.Md5(fmt.Sprintf("%s%v%d%s", config.Configs.Server.Name, id, now.Timestamp(), helper.StringRandom(8)))
}
