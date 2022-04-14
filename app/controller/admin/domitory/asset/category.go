package asset

import (
	"github.com/gin-gonic/gin"
	"saas/app/constant"
	"saas/app/form/admin/dormitory/asset"
	"saas/app/model"
	assetResponse "saas/app/response/admin/dormitory/asset"
	"saas/kernel/data"
	"saas/kernel/response"
	"strconv"
)

func DoCategoryByCreate(ctx *gin.Context) {

	var former asset.DoCategoryByCreateForm
	if err := ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	category := model.DorAssetCategory{
		Name:     former.Name,
		Order:    former.Order,
		IsEnable: former.IsEnable,
	}

	if data.Database.Create(&category); category.Id <= 0 {
		response.ToResponseByFail(ctx, "添加失败")
		return
	}

	response.ToResponseBySuccess(ctx)
}

func DoCategoryByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.ToResponseByFailRequestMessage(ctx, "ID获取失败")
		return
	}

	var former asset.DoCategoryByUpdateForm
	if err = ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var category model.DorAssetCategory
	data.Database.First(&category, id)
	if category.Id <= 0 {
		response.ToResponseByNotFound(ctx, "未找到该类型")
		return
	}

	category.Name = former.Name
	category.Order = former.Order
	category.IsEnable = former.IsEnable

	if t := data.Database.Save(&category); t.RowsAffected <= 0 {
		response.ToResponseByFail(ctx, "修改失败")
		return
	}

	response.ToResponseBySuccess(ctx)
}

func DoCategoryByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.ToResponseByFailRequestMessage(ctx, "ID获取失败")
		return
	}

	var category model.DorAssetCategory
	data.Database.First(&category, id)
	if category.Id <= 0 {
		response.ToResponseByNotFound(ctx, "未找到该类型")
		return
	}

	if t := data.Database.Delete(&category); t.RowsAffected <= 0 {
		response.ToResponseByFail(ctx, "删除失败")
		return
	}

	response.ToResponseBySuccess(ctx)
}

func DoCategoryByEnable(ctx *gin.Context) {

	var former asset.DoCategoryByEnableForm
	if err := ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var category model.DorAssetCategory
	data.Database.First(&category, former.Id)
	if category.Id <= 0 {
		response.ToResponseByNotFound(ctx, "未找到该类型")
		return
	}

	category.IsEnable = former.IsEnable

	if t := data.Database.Save(&category); t.RowsAffected <= 0 {
		response.ToResponseByFail(ctx, "启禁失败")
		return
	}

	response.ToResponseBySuccess(ctx)
}

func ToCategoryByList(ctx *gin.Context) {

	responses := make([]any, 0)

	var categories []model.DorAssetCategory
	data.Database.Order("`order` asc, `id` desc").Find(&categories)

	for _, item := range categories {
		responses = append(responses, assetResponse.ToCategoryByListResponse{
			Id:        item.Id,
			Name:      item.Name,
			Order:     item.Order,
			IsEnable:  item.IsEnable,
			CreatedAt: item.CreatedAt.ToDateTimeString(),
		})
	}

	response.ToResponseBySuccessList(ctx, responses)
}

func ToCategoryByOnline(ctx *gin.Context) {

	responses := make([]any, 0)

	var categories []model.DorAssetCategory
	data.Database.Where("`is_enable`=?", constant.IsEnableYes).Order("`order` asc, `id` desc").Find(&categories)

	for _, item := range categories {
		responses = append(responses, assetResponse.ToCategoryByOnlineResponse{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	response.ToResponseBySuccessList(ctx, responses)
}
