package kernel

import (
	"saas/app/crontab"
	"saas/kernel/api"
	"saas/kernel/authorize"
	"saas/kernel/cmd"
	"saas/kernel/config"
	"saas/kernel/database"
	"saas/kernel/logger"
	"saas/kernel/snowflake"
)

var services = []func(){
	config.InitConfig,
	api.InitApi,
	logger.InitLogger,
	database.InitMySQL,
	database.InitRedis,
	authorize.InitCasbin,
	snowflake.InitSnowflake,
	crontab.InitCrontab,
	cmd.InitCmd,
}

func Bootstrap() {

	for _, item := range services {
		item()
	}

}
