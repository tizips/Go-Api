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
	"saas/kernel/validator"
	"saas/routes"
	"time"
)

func Server() {

	gin.SetMode(config.Values.Server.Mode)

	fmt.Printf("Listen: %d\n", config.Values.Server.Port)

	app := gin.New()

	validator.Init()

	routes.Routes(app)

	config.Application.Application = app

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Values.Server.Port),
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
