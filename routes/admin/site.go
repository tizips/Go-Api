package admin

import (
	"github.com/gin-gonic/gin"
	"saas/app/controller/admin/site/architecture"
	"saas/app/controller/admin/site/helper"
	"saas/app/controller/admin/site/manage"
)

func RouteSite(route *gin.RouterGroup) {

	site := route.Group("site")
	{
		helperGroup := site.Group("helper")
		{
			helperGroup.GET("apis", helper.ToApiByList)
		}

		manageGroup := site.Group("manage")
		{
			admins := manageGroup.Group("admins")
			{
				admins.GET("", manage.ToAdminByPaginate)
				admins.PUT(":id", manage.DoAdminByUpdate)
				admins.DELETE(":id", manage.DoAdminByDelete)
			}

			admin := manageGroup.Group("admin")
			{
				admin.POST("", manage.DoAdminByCreate)
				admin.PUT("enable", manage.DoAdminByEnable)
			}

			permissions := manageGroup.Group("permissions")
			{
				permissions.GET("", manage.ToPermissionByTree)
				permissions.PUT(":id", manage.DoPermissionByUpdate)
				permissions.DELETE(":id", manage.DoPermissionByDelete)
			}

			permission := manageGroup.Group("permission")
			{
				permission.GET("parents", manage.ToPermissionByParents)
				permission.GET("self", manage.ToPermissionBySelf)
				permission.POST("", manage.DoPermissionByCreate)
			}

			roles := manageGroup.Group("roles")
			{
				roles.GET("", manage.ToRoleByPaginate)
				roles.PUT(":id", manage.DoRoleByUpdate)
				roles.DELETE(":id", manage.DoRoleByDelete)
			}

			role := manageGroup.Group("role")
			{
				role.POST("", manage.DoRoleByCreate)
				role.GET("online", manage.ToRoleByOnline)
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
