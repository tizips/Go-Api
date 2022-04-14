package routes

import (
	"github.com/gin-gonic/gin"
	"saas/app/middleware/basic"
	"saas/routes/admin"
)

func Routes(route *gin.Engine) {

	route.Use(basic.LoggerMiddleware())

	admin.Admins(route)
}
