package basic

import (
	"github.com/gin-gonic/gin"
	"saas/app/request/admin/basic"
	basicResponse "saas/app/response/admin/basic"
	"saas/app/service/helper"
	"saas/kernel/response"
)

func DoUploadBySimple(ctx *gin.Context) {

	var request basic.DoUploadBySimple

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	file, err := ctx.FormFile("file")

	if err != nil {
		response.Fail(ctx, "上传失败，请稍后重试")
		return
	}

	result, err := helper.DoUploadBySimple(ctx, "/system/"+request.Dir, file)

	if err != nil || result == nil {
		response.Fail(ctx, "上传失败，请稍后重试")
		return
	}

	response.SuccessByData(ctx, basicResponse.DoUploadBySimple{
		Name: result.Name,
		Path: result.Path,
		Url:  result.Url,
	})
}
