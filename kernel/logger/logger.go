package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"saas/kernel/config"
	"time"
)

var Logger struct {
	Api       *logrus.Logger
	SQL       *logrus.Logger
	Exception *logrus.Logger
	Amqp      *logrus.Logger
}

func InitLogger() {

	folder()

	api()

	sql()

	exception()

	amqp()

}

func folder() {

	path := path("")

	if err := os.MkdirAll(path, 0750); err != nil {
		fmt.Printf("日志文件夹创建失败:%s\nerror:%v", path, err)
		os.Exit(1)
	}

}

func api() {

	filename := path("api")

	if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
		if _, err := os.Create(filename); err != nil {
			fmt.Printf("日志文件创建失败:%s\nerror:%v", filename, err)
			os.Exit(1)
		}
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

	if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
		if _, err := os.Create(filename); err != nil {
			fmt.Printf("日志文件创建失败:%s\nerror:%v", filename, err)
			os.Exit(1)
		}
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

func exception() {

	filename := path("exception")

	if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
		if _, err := os.Create(filename); err != nil {
			fmt.Printf("日志文件创建失败:%s\nerror:%v", filename, err)
			os.Exit(1)
		}
	}

	writer, _ := rotatelogs.New(
		filename+".%Y%m%d",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*30),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	Logger.Exception = logrus.New()

	Logger.Exception.SetFormatter(&logrus.JSONFormatter{})
	Logger.Exception.Hooks = make(logrus.LevelHooks)
	Logger.Exception.ExitFunc = os.Exit
	Logger.Exception.SetOutput(writer)

}

func amqp() {

	filename := path("amqp")

	if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
		if _, err := os.Create(filename); err != nil {
			fmt.Printf("日志文件创建失败:%s\nerror:%v", filename, err)
			os.Exit(1)
		}
	}

	writer, _ := rotatelogs.New(
		filename+".%Y%m%d",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*30),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	Logger.Amqp = logrus.New()

	Logger.Amqp.SetFormatter(&logrus.JSONFormatter{})
	Logger.Amqp.Hooks = make(logrus.LevelHooks)
	Logger.Amqp.ExitFunc = os.Exit
	Logger.Amqp.SetOutput(writer)

}

func path(filename string) string {

	filepath := fmt.Sprintf("%s/logs", config.Application.Runtime)

	if filename != "" {
		filepath += "/" + filename + ".log"
	}

	return filepath
}
