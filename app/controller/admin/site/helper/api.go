package helper

import (
	"github.com/gin-gonic/gin"
	"saas/app/model"
	"saas/app/request/admin/site/helper"
	res "saas/app/response/admin/site/helper"
	"saas/kernel/api"
	"saas/kernel/app"
	"saas/kernel/response"
	"strings"
)

func ToApiByList(ctx *gin.Context) {

	var request helper.ToApiByList

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var module model.SysModule

	if app.MySQL.Find(&module, request.Module); module.Id <= 0 {
		response.NotFound(ctx, "模块不存在")
		return
	}

	routes := app.Engine.Routes()

	var responses = make([]res.ToApiByList, 0)

	var permissions []model.SysPermission

	app.MySQL.Find(&permissions, "`module_id`=? and `method`<>? and `path`<>?", module.Id, "", "")

	var permissionsCache = make(map[string]bool, 0)

	if len(permissions) > 0 {
		for _, item := range permissions {
			permissionsCache[api.OmitKey(item.Method, item.Path)] = true
		}
	}

	for _, item := range routes {

		urls := strings.Split(item.Path, "/")

		if len(urls) >= 3 && urls[1] == "admin" && urls[2] == module.Slug {
			mark := true
			if _, exist := api.OmitOfCache[api.OmitKey(item.Method, item.Path)]; exist {
				mark = false
			}

			if mark && len(permissionsCache) > 0 {
				if _, exist := permissionsCache[api.OmitKey(item.Method, item.Path)]; exist {
					mark = false
				}
			}

			if mark {
				responses = append(responses, res.ToApiByList{
					Method: item.Method,
					Path:   item.Path,
				})
			}
		}
	}

	response.SuccessByData(ctx, responses)
}
