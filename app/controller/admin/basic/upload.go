package basic

import (
	"github.com/gin-gonic/gin"
	"saas/app/request/admin/basic"
	res "saas/app/response/admin/basic"
	"saas/kernel/response"
	"saas/kit/filesystem"
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

	storage := filesystem.New().Upload()

	uri, name, err := storage.Save(file, "/system/"+request.Dir, "")

	if err != nil {
		response.Fail(ctx, "上传失败，请稍后重试")
		return
	}

	responses := res.DoUploadBySimple{
		Name: name,
		Path: uri,
		Url:  storage.Url(uri),
	}

	response.SuccessByData(ctx, responses)
}
