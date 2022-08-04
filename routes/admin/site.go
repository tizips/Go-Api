package admin

import (
	"github.com/gin-gonic/gin"
	"saas/app/controller/admin/site/architecture"
	"saas/app/controller/admin/site/helper"
	"saas/app/controller/admin/site/management"
)

func RouteSite(route *gin.RouterGroup) {

	site := route.Group("site")
	{
		helperGroup := site.Group("helper")
		{
			helperGroup.GET("apis", helper.ToApiByList)
		}

		managementGroup := site.Group("management")
		{
			admins := managementGroup.Group("admins")
			{
				admins.GET("", management.ToAdminByPaginate)
				admins.PUT(":id", management.DoAdminByUpdate)
				admins.DELETE(":id", management.DoAdminByDelete)
			}

			admin := managementGroup.Group("admin")
			{
				admin.POST("", management.DoAdminByCreate)
				admin.PUT("enable", management.DoAdminByEnable)
			}

			permissions := managementGroup.Group("permissions")
			{
				permissions.GET("", management.ToPermissionByTree)
				permissions.PUT(":id", management.DoPermissionByUpdate)
				permissions.DELETE(":id", management.DoPermissionByDelete)
			}

			permission := managementGroup.Group("permission")
			{
				permission.GET("parents", management.ToPermissionByParents)
				permission.GET("self", management.ToPermissionBySelf)
				permission.POST("", management.DoPermissionByCreate)
			}

			roles := managementGroup.Group("roles")
			{
				roles.GET("", management.ToRoleByPaginate)
				roles.PUT(":id", management.DoRoleByUpdate)
				roles.DELETE(":id", management.DoRoleByDelete)
			}

			role := managementGroup.Group("role")
			{
				role.POST("", management.DoRoleByCreate)
				role.GET("enable", management.ToRoleByEnable)
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
