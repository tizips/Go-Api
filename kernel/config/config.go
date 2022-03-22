package config

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"saas/app/helper/str"
	"strings"
)

var Configs struct {
	System   System
	Database Database
	Server   Server
	Redis    Redis
	Jwt      Jwt
	Cache    Cache
}

type System struct {
	Application *gin.Engine
	Path        string
	Public      string
}

func InitConfig() {

	pwd, _ := os.Getwd()

	Configs.System = System{
		Path:   pwd,
		Public: pwd + "/public",
	}

	cfg, err := ini.Load("conf/env.conf")

	if err != nil {
		fmt.Printf("Fail to load env file: %v", err)
		os.Exit(1)
	}

	files, err := ioutil.ReadDir(Configs.System.Path + "/kernel/config")

	if err != nil {
		fmt.Printf("Fail to load config file: %v", err)
		os.Exit(1)
	}

	for _, item := range files {
		if !item.IsDir() {
			filename := strings.TrimSuffix(path.Base(item.Name()), path.Ext(item.Name()))

			if filename != "config" {

				v := reflect.ValueOf(&Configs).Elem()

				for i := 0; i < v.NumField(); i++ {

					if strings.Title(filename) == v.Field(i).Type().Name() {

						for j := 0; j < v.Field(i).Type().NumField(); j++ {

							defaults := v.Field(i).Type().Field(j).Tag.Get("default")
							name := str.Snake(v.Field(i).Type().Field(j).Name)
							value := cfg.Section(filename).Key(name).Value()

							if value == "" && defaults != "omitempty" {
								v.Field(i).Field(j).Set(reflect.ValueOf(defaults))
							} else if value != "" {
								v.Field(i).Field(j).Set(reflect.ValueOf(value))
							}
						}
					}
				}
			}
		}
	}
}
