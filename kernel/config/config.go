package config

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
	"os"
	"reflect"
	"saas/app/helper/str"
	"saas/kernel/config/configs"
	"strconv"
)

var Application system

var Values struct {
	Server   configs.Server
	Database configs.Database
	Redis    configs.Redis
	Jwt      configs.Jwt
	Cache    configs.Cache
}

type system struct {
	Application *gin.Engine
	Path        string
	Public      string
}

func InitConfig() {

	pwd, _ := os.Getwd()

	Application = system{
		Path:   pwd,
		Public: pwd + "/public",
	}

	handler()

}

func handler() {

	cfg, err := ini.Load(Application.Path + "/conf/env.conf")

	if err != nil {
		fmt.Printf("Fail to load env file: %v", err)
		os.Exit(1)
	}

	e := reflect.ValueOf(&Values).Elem()

	for i := 0; i < e.NumField(); i++ {

		section := str.Snake(e.Field(i).Type().Name())

		for j := 0; j < e.Field(i).Type().NumField(); j++ {

			key := str.Snake(e.Field(i).Type().Field(j).Name)

			defaults := e.Field(i).Type().Field(j).Tag.Get("default")

			types := e.Field(i).Type().Field(j).Type.Kind()

			values := cfg.Section(section).Key(key)

			switch types {
			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
				val, _ := values.Int64()
				if val <= 0 {
					val, _ = strconv.ParseInt(defaults, 10, 64)
				}
				e.Field(i).Field(j).SetInt(val)
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
				val, _ := values.Uint64()
				if val <= 0 {
					val, _ = strconv.ParseUint(defaults, 10, 64)
				}
				e.Field(i).Field(j).SetUint(val)
			case reflect.Float32, reflect.Float64:
				val, _ := values.Float64()
				if val <= 0 {
					val, _ = strconv.ParseFloat(defaults, 32)
				}
				e.Field(i).Field(j).SetFloat(val)
			case reflect.Bool:
				val := values.MustBool(defaults == "true")
				e.Field(i).Field(j).SetBool(val)
			case reflect.String:
				val := values.String()
				if val == "" && defaults != "" {
					val = defaults
				}
				e.Field(i).Field(j).SetString(val)
			}
		}
	}
}
