package routes

import (
	"github.com/gin-gonic/gin"
	"saas/app/middleware/basic"
	"saas/kernel/app"
	"saas/routes/admin"
)

func Routes(route *gin.Engine) {

	route.Use(basic.LoggerMiddleware())

	route.Static("/upload", app.Dir.Runtime+"/upload")

	admin.Admins(route)
}
