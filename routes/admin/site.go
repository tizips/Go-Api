package admin

import (
	"github.com/gin-gonic/gin"
	"saas/app/controller/admin/site/architecture"
	"saas/app/controller/admin/site/auth"
	"saas/app/controller/admin/site/helper"
)

func RouteSite(route *gin.RouterGroup) {

	site := route.Group("site")
	{
		helperGroup := site.Group("helper")
		{
			helperGroup.GET("apis", helper.ToApiByList)
		}

		authGroup := site.Group("auth")
		{
			admins := authGroup.Group("admins")
			{
				admins.GET("", auth.ToAdminByPaginate)
				admins.PUT(":id", auth.DoAdminByUpdate)
				admins.DELETE(":id", auth.DoAdminByDelete)
			}

			admin := authGroup.Group("admin")
			{
				admin.POST("", auth.DoAdminByCreate)
				admin.PUT("enable", auth.DoAdminByEnable)
			}

			permissions := authGroup.Group("permissions")
			{
				permissions.GET("", auth.ToPermissionByTree)
				permissions.PUT(":id", auth.DoPermissionByUpdate)
				permissions.DELETE(":id", auth.DoPermissionByDelete)
			}

			permission := authGroup.Group("permission")
			{
				permission.GET("parents", auth.ToPermissionByParents)
				permission.GET("self", auth.ToPermissionBySelf)
				permission.POST("", auth.DoPermissionByCreate)
			}

			roles := authGroup.Group("roles")
			{
				roles.GET("", auth.ToRoleByPaginate)
				roles.PUT(":id", auth.DoRoleByUpdate)
				roles.DELETE(":id", auth.DoRoleByDelete)
			}

			role := authGroup.Group("role")
			{
				role.POST("", auth.DoRoleByCreate)
				role.GET("enable", auth.ToRoleByEnable)
			}
		}

		architecturesGroup := site.Group("architecture")
		{
			modules := architecturesGroup.Group("modules")
			{
				modules.GET("", architecture.ToModuleByList)
				modules.PUT(":id", architecture.DoModuleByUpdate)
				modules.DELETE(":id", architecture.DoModuleByDelete)
			}

			module := architecturesGroup.Group("module")
			{
				module.POST("", architecture.DoModuleByCreate)
				module.GET("online", architecture.ToModuleByOnline)
				module.PUT("enable", architecture.DoModuleByEnable)
			}
		}
	}
}
