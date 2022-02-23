package routes

import (
	"github.com/gin-gonic/gin"
	"saas/app/controller/admin/account"
	"saas/app/middleware/basic"
	adminModule "saas/routes/admin"
)

func Admins(route *gin.Engine) {
	admin := route.Group("/admin")
	admin.Use(basic.JwtParse())
	{
		login := admin.Group("/login")
		{
			login.POST("/account", account.DoLoginByAccount)
			login.POST("/qrcode", account.DoLoginByQrcode)
		}

		auth := admin.Group("")
		auth.Use(basic.Auth("admin"))
		{
			accountGroup := auth.Group("/account")
			{
				accountGroup.GET("/information", account.ToAccountByInformation)
				accountGroup.GET("/module", account.ToAccountByModule)
				accountGroup.GET("/permission", account.ToAccountByPermission)
			}

			adminModule.RouteSite(auth)
			adminModule.RouteDormitory(auth)
		}
	}
}
