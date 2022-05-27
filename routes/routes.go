package routes

import (
	"github.com/gin-gonic/gin"
	"saas/app/middleware/basic"
	"saas/kernel/config"
	"saas/routes/admin"
)

func Routes(route *gin.Engine) {

	route.Use(basic.LoggerMiddleware())

	route.Static("/upload", config.Application.Runtime+"/upload")

	admin.Admins(route)
}
