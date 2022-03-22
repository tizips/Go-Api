package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"saas/kernel/config"
	"saas/kernel/dir"
	"time"
)

var Logger struct {
	Api *logrus.Logger
	SQL *logrus.Logger
}

func InitLogger() {

	fmt.Println(config.Configs.System.Path)

	folder()

	api()

	sql()

}

func folder() {

	path := path("")

	err := dir.Mkdir(path)
	if err != nil {
		fmt.Printf("日志文件夹创建失败:%s\nerror:%v", path, err)
		os.Exit(1)
	}

}

func api() {

	filename := path("api")

	_, err := dir.Touch(filename)
	if err != nil {
		fmt.Printf("日志文件创建失败:%s\nerror:%v", filename, err)
		os.Exit(1)
	}

	writer, _ := rotatelogs.New(
		filename+".%Y%m%d",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*30),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	Logger.Api = logrus.New()

	Logger.Api.SetFormatter(&logrus.JSONFormatter{})
	Logger.Api.SetOutput(writer)
}

func sql() {

	filename := path("sql")

	_, err := dir.Touch(filename)
	if err != nil {
		fmt.Printf("日志文件创建失败:%s\nerror:%v", filename, err)
		os.Exit(1)
	}

	writer, _ := rotatelogs.New(
		filename+".%Y%m%d",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*30),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	Logger.SQL = logrus.New()

	Logger.SQL.SetFormatter(&logrus.JSONFormatter{})
	Logger.SQL.Hooks = make(logrus.LevelHooks)
	Logger.SQL.ExitFunc = os.Exit
	Logger.SQL.SetOutput(writer)

}

func path(filename string) string {

	path := fmt.Sprintf("%s/logs", config.Configs.System.Public)
	if filename != "" {
		path += "/" + filename + ".log"
	}

	return path
}
