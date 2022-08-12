package kernel

import (
	"github.com/spf13/cobra"
	"saas/app/crontab"
	"saas/kernel/api"
	"saas/kernel/authorize"
	"saas/kernel/config"
	"saas/kernel/data"
	"saas/kernel/logger"
	"saas/kernel/snowflake"
	"saas/kernel/system"
)

func Application() {

	initialize()

	cmd := &cobra.Command{
		Use:     "saas",
		Short:   "一站式应用解决方案",
		Version: "1.0.0",
	}

	server := &cobra.Command{
		Use:   "server",
		Short: "运行应用程序",
		Run: func(cmd *cobra.Command, args []string) {
			system.Server()
		},
	}

	root := &cobra.Command{
		Use:   "root",
		Short: "生成系统开发账号",
		Run: func(cmd *cobra.Command, args []string) {
			system.Root()
		},
	}

	password := &cobra.Command{
		Use:   "password",
		Short: "生成密码",
		Run: func(cmd *cobra.Command, args []string) {
			system.Password()
		},
	}

	cmd.AddCommand(server)
	cmd.AddCommand(root)
	cmd.AddCommand(password)

	_ = cmd.Execute()

}

func initialize() {

	config.InitConfig()

	api.InitApi()

	logger.InitLogger()

	data.InitDatabase()

	data.InitRedis()

	authorize.InitCasbin()

	_ = snowflake.InitSnowflake()

	go crontab.InitCrontab()

}
