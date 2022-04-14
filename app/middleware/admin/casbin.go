package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"saas/kernel/api"
	"saas/kernel/auth"
	"saas/kernel/response"
)

func CasbinMiddleware() gin.HandlerFunc {

	return func(ctx *gin.Context) {

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
	_, exist := api.OmitOfCache[api.OmitKey(method, path)]
	return exist
}
