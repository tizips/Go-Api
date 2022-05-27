package helper

import (
	"github.com/gin-gonic/gin"
	"saas/app/model"
	helperForm "saas/app/request/admin/site/helper"
	"saas/app/response/admin/site/helper"
	"saas/kernel/api"
	"saas/kernel/config"
	"saas/kernel/data"
	"saas/kernel/response"
	"strings"
)

func ToApiByList(ctx *gin.Context) {

	var request helperForm.ToApiByList
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var module model.SysModule
	data.Database.Find(&module, request.Module)
	if module.Id <= 0 {
		response.NotFound(ctx, "模块不存在")
		return
	}

	routes := config.Application.Application.Routes()

	var responses = make([]any, 0)

	var permissions []model.SysPermission
	data.Database.
		Where("`module_id` = ? and `method` <> ? and `path` <> ?", module.Id, "", "").
		Find(&permissions)

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
				responses = append(responses, helper.ToApiByList{
					Method: item.Method,
					Path:   item.Path,
				})
			}
		}
	}

	response.SuccessByList(ctx, responses)
}
