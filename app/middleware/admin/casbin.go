package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"saas/app/controller/admin/site/helper"
	"saas/kernel/auth"
	"saas/kernel/response"
)

func CasbinMiddleware() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		// @todo 缓存被忽略的路由，直接从内存中查询
		//	判断该路由是否被忽略权限判断中
		if !omit(ctx.Request.Method, ctx.FullPath()) {

			//	判断该用户是否为开发组权限
			if ok, _ := auth.Casbin.HasRoleForUser(auth.NameByAdmin(auth.Id(ctx)), auth.NameByRole(auth.ROOT)); !ok {

				//	判断该用户是否有该接口的访问权限
				if ok, _ = auth.Casbin.Enforce(auth.NameByAdmin(auth.Id(ctx)), ctx.Request.Method, ctx.FullPath()); !ok {
					ctx.Abort()
					ctx.JSON(http.StatusForbidden, response.Response{
						Code:    40300,
						Message: "Forbidden",
					})
					return
				}
			}
		}

		ctx.Next()
	}
}

func omit(method string, path string) bool {

	for _, item := range helper.Omits() {
		if method == item.Method && path == item.Path {
			return true
		}
	}

	return false
}
