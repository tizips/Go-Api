package admin

import (
	"github.com/gin-gonic/gin"
	"saas/kernel/api"
	"saas/kernel/app"
	"saas/kernel/authorize"
	"saas/kernel/response"
)

func CasbinMiddleware() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		//	判断该路由是否被忽略权限判断中
		if !omit(ctx.Request.Method, ctx.FullPath()) {

			//	判断该用户是否为开发组权限
			if !authorize.Root(authorize.Id(ctx)) {

				//	判断该用户是否有该接口的访问权限
				if ok, _ := app.Casbin.Enforce(authorize.NameByAdmin(authorize.Id(ctx)), ctx.Request.Method, ctx.FullPath()); !ok {
					ctx.Abort()
					response.Forbidden(ctx)
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
