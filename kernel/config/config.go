package config

import (
	"github.com/creasty/defaults"
	"github.com/gookit/color"
	"gopkg.in/yaml.v3"
	"os"
	"saas/kernel/app"
)

func InitConfig() {

	pwd, _ := os.Getwd()

	app.Dir.Root = pwd
	app.Dir.Runtime = pwd + "/runtime"

	file, err := os.ReadFile(app.Dir.Root + "/conf/env.yaml")
	if err != nil {
		color.Errorf("Fail to load env file: %v\n", err)
		os.Exit(1)
	}

	if err := yaml.Unmarshal(file, &app.Cfg); err != nil {
		color.Errorf("Fail to parse env file: %v\n", err)
		os.Exit(1)
	}

	if err := defaults.Set(&app.Cfg); err != nil {
		color.Errorf("Fail to default env params: %v\n", err)
		os.Exit(1)
	}
}
