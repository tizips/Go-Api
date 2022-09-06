package logger

import (
	"fmt"
	"github.com/gookit/color"
	rotate "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"saas/kernel/app"
	"saas/kit/filesystem"
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

	storage := filesystem.Disk("local")

	if err := storage.Mkdir("logs"); err != nil {
		color.Errorln("fail to mkdir logs: %v\n", err)
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

	writer, _ := rotate.New(
		filename+".%Y%m%d",
		rotate.WithLinkName(filename),
		rotate.WithMaxAge(time.Hour*24*30),
		rotate.WithRotationTime(time.Hour*24),
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

	writer, _ := rotate.New(
		filename+".%Y%m%d",
		rotate.WithLinkName(filename),
		rotate.WithMaxAge(time.Hour*24*30),
		rotate.WithRotationTime(time.Hour*24),
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

	writer, _ := rotate.New(
		filename+".%Y%m%d",
		rotate.WithLinkName(filename),
		rotate.WithMaxAge(time.Hour*24*30),
		rotate.WithRotationTime(time.Hour*24),
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

	writer, _ := rotate.New(
		filename+".%Y%m%d",
		rotate.WithLinkName(filename),
		rotate.WithMaxAge(time.Hour*24*30),
		rotate.WithRotationTime(time.Hour*24),
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
