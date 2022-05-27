package kernel

import (
	"os"
	"saas/app/amqp"
	"saas/app/crontab"
	"saas/kernel/api"
	"saas/kernel/authorize"
	"saas/kernel/config"
	"saas/kernel/data"
	"saas/kernel/dir"
	"saas/kernel/logger"
	"saas/kernel/system"
)

func Application() {

	servers := os.Args[1:]

	if len(servers) <= 0 {
		system.Help()
		return
	}

	initialize()

	switch servers[0] {
	case "server":
		system.Server()
	case "root":
		system.Root()
	case "password":
		system.Password()
	default:
		system.Help()
	}

}

func initialize() {

	config.InitConfig()

	api.InitApi()

	dir.InitDir()

	logger.InitLogger()

	data.InitDatabase()

	data.InitRedis()

	authorize.InitCasbin()

	go crontab.InitCrontab()

	go amqp.Init()

}
