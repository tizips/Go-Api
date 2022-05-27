package admin

import (
	"github.com/gin-gonic/gin"
	"saas/app/controller/admin/basic"
	adminMiddleware "saas/app/middleware/admin"
	basicMiddleware "saas/app/middleware/basic"
)

func Admins(route *gin.Engine) {
	admin := route.Group("/admin")
	admin.Use(basicMiddleware.JwtParseMiddleware())
	{
		login := admin.Group("/login")
		{
			login.POST("/account", basicMiddleware.LimitMiddleware(nil), basic.DoLoginByAccount)
			login.POST("/qrcode", basicMiddleware.LimitMiddleware(nil), basic.DoLoginByQrcode)
		}

		auth := admin.Group("")
		auth.Use(basicMiddleware.AuthMiddleware(), adminMiddleware.CasbinMiddleware())
		{
			accountGroup := auth.Group("/account")
			{
				accountGroup.PUT("", basic.DoAccountByUpdate)
				accountGroup.GET("/information", basic.ToAccountByInformation)
				accountGroup.GET("/module", basic.ToAccountByModule)
				accountGroup.GET("/permission", basic.ToAccountByPermission)
				accountGroup.POST("/logout", basic.DoLogout)
			}

			RouteSite(auth)
			RouteDormitory(auth)

			uploadGroup := auth.Group("upload")
			{
				uploadGroup.POST("", basic.DoUploadBySimple)
			}
		}
	}
}
