package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"saas/kernel/app"
	"time"
)

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

	app.Logger.Api = logrus.New()

	app.Logger.Api.SetFormatter(&logrus.JSONFormatter{})
	app.Logger.Api.SetOutput(writer)
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

	app.Logger.SQL = logrus.New()

	app.Logger.SQL.SetFormatter(&logrus.JSONFormatter{})
	app.Logger.SQL.Hooks = make(logrus.LevelHooks)
	app.Logger.SQL.ExitFunc = os.Exit
	app.Logger.SQL.SetOutput(writer)

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

	app.Logger.Exception = logrus.New()

	app.Logger.Exception.SetFormatter(&logrus.JSONFormatter{})
	app.Logger.Exception.Hooks = make(logrus.LevelHooks)
	app.Logger.Exception.ExitFunc = os.Exit
	app.Logger.Exception.SetOutput(writer)

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

	app.Logger.Amqp = logrus.New()

	app.Logger.Amqp.SetFormatter(&logrus.JSONFormatter{})
	app.Logger.Amqp.Hooks = make(logrus.LevelHooks)
	app.Logger.Amqp.ExitFunc = os.Exit
	app.Logger.Amqp.SetOutput(writer)

}

func path(filename string) string {

	filepath := fmt.Sprintf("%s/logs", app.Dir.Runtime)

	if filename != "" {
		filepath += "/" + filename + ".log"
	}

	return filepath
}
