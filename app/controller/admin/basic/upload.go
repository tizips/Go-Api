package basic

import (
	"github.com/gin-gonic/gin"
	"saas/app/request/admin/basic"
	basicResponse "saas/app/response/admin/basic"
	"saas/app/service/helper"
	basic3 "saas/kernel/response"
)

func DoUploadBySimple(ctx *gin.Context) {

	var request basic.DoUploadBySimple
	if err := ctx.ShouldBind(&request); err != nil {
		basic3.FailByRequest(ctx, err)
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		basic3.Fail(ctx, "上传失败，请稍后重试")
		return
	}

	err, result := helper.DoUploadBySimple(ctx, "/system/"+request.Dir, file)

	if err != nil {
		basic3.Fail(ctx, "上传失败，请稍后重试")
		return
	}

	basic3.SuccessByData(ctx, basicResponse.DoUploadBySimple{
		Name: result.Name,
		Path: result.Path,
		Url:  result.Url,
	})
}
