package helper

import (
	"github.com/gin-gonic/gin"
	"net/http"
	helperForm "saas/app/form/admin/site/helper"
	"saas/app/model"
	"saas/app/response/admin/site/helper"
	"saas/kernel/config"
	"saas/kernel/data"
	"saas/kernel/response"
	"strings"
)

func ToApiByList(ctx *gin.Context) {

	var former helperForm.ToApiByListForm
	if err := ctx.ShouldBindQuery(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var module model.SysModule
	data.Database.First(&module, former.Module)
	if module.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "模块不存在",
		})
		return
	}

	routes := config.Configs.System.Application.Routes()

	res := response.Responses{
		Code:    20000,
		Message: "Success",
		Data:    []any{},
	}

	var permissions []model.SysPermission
	data.Database.
		Where("module_id = ?", module.Id).
		Where("method <> ?", "").
		Where("path <> ?", "").
		Find(&permissions)

	for _, item := range routes {
		urls := strings.Split(item.Path, "/")
		if len(urls) >= 3 && urls[1] == "admin" && urls[2] == module.Slug {
			mark := true
			for _, value := range Omits() {
				if item.Method == value.Method && item.Path == value.Path {
					mark = false
				}
			}
			if mark {
				for _, value := range permissions {
					if item.Method == value.Method && item.Path == value.Path {
						mark = false
					}
				}
			}
			if mark {
				res.Data = append(res.Data, helper.ToApiByListResponse{
					Method: item.Method,
					Path:   item.Path,
				})
			}
		}
	}

	ctx.JSON(http.StatusOK, res)
}

//	被忽略返回的 Api
func Omits() []helper.ToApiByListResponse {
	return []helper.ToApiByListResponse{

		{Method: http.MethodPost, Path: "/admin/login/account"},
		{Method: http.MethodPost, Path: "/admin/login/qrcode"},
		{Method: http.MethodGet, Path: "/admin/account/information"},
		{Method: http.MethodGet, Path: "/admin/account/module"},
		{Method: http.MethodGet, Path: "/admin/account/permission"},

		{Method: http.MethodGet, Path: "/admin/site/helper/apis"},
		{Method: http.MethodGet, Path: "/admin/site/auth/permission/parents"},
		{Method: http.MethodGet, Path: "/admin/site/auth/permission/self"},
		{Method: http.MethodGet, Path: "/admin/site/auth/role/enable"},
		{Method: http.MethodGet, Path: "/admin/site/architecture/module/online"},

		{Method: http.MethodGet, Path: "/admin/dormitory/basic/type/online"},
		{Method: http.MethodGet, Path: "/admin/dormitory/basic/building/online"},
		{Method: http.MethodGet, Path: "/admin/dormitory/basic/floor/online"},
		{Method: http.MethodGet, Path: "/admin/dormitory/basic/room/online"},
		{Method: http.MethodGet, Path: "/admin/dormitory/basic/bed/online"},
		{Method: http.MethodGet, Path: "/admin/dormitory/stay/category/online"},
		{Method: http.MethodGet, Path: "/admin/dormitory/asset/category/online"},
		{Method: http.MethodGet, Path: "/admin/dormitory/asset/device/online"},
		{Method: http.MethodGet, Path: "/admin/dormitory/asset/package/online"},
	}
}
