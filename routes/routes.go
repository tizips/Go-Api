package routes

import (
	"github.com/gin-gonic/gin"
	"saas/app/middleware/basic"
)

func Routes(route *gin.Engine) {

	route.Use(basic.LoggerMiddleware())

	Admins(route)
}
