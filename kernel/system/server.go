package system

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"saas/kernel/config"
	"saas/routes"
	"time"
)

func Server() {

	gin.SetMode(config.Configs.Server.Mode)

	fmt.Printf("Listen: %s\n", config.Configs.Server.Port)

	app := gin.New()

	routes.Routes(app)

	config.Configs.System.Application = app

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Configs.Server.Port),
		Handler: app,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
