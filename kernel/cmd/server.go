package cmd

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/signal"
	"saas/kernel/app"
	"saas/kernel/validator"
	"saas/routes"
	"time"
)

func Server(command *cobra.Command) {

	cmd := &cobra.Command{
		Use:   "server",
		Short: "运行应用程序",
		Run:   server,
	}

	command.AddCommand(cmd)
}

func server(cmd *cobra.Command, args []string) {

	gin.SetMode(app.Cfg.Server.Mode)

	engine := gin.New()

	validator.Init()

	routes.Routes(engine)

	app.Engine = engine

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Cfg.Server.Port),
		Handler: engine,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			color.Warnf("listen: %s\n", err)
		}
	}()

	color.Warnf("\nApplication Listen: %d\n\n", app.Cfg.Server.Port)

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	color.Warnln("\nShutdown Server ...\n")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		color.Errorf("\nServer Shutdown:", err)
	}

	color.Errorln("\nServer exiting")
}
