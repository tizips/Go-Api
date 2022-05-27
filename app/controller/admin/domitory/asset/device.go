package asset

import (
	"github.com/gin-gonic/gin"
	"saas/app/constant"
	"saas/app/model"
	"saas/app/request/admin/dormitory/asset"
	asset2 "saas/app/response/admin/dormitory/asset"
	"saas/kernel/data"
	"saas/kernel/response"
	"strconv"
)

func DoDeviceByCreate(ctx *gin.Context) {

	var request asset.DoDeviceByCreate
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var count int64 = 0
	data.Database.Model(model.DorAssetCategory{}).Where("id=?", request.Category).Where("is_enable=?", constant.IsEnableYes).Count(&count)

	if count <= 0 {
		response.NotFound(ctx, "类型不存在")
		return
	}

	assets := model.DorDevice{
		CategoryId:    request.Category,
		No:            request.No,
		Name:          request.Name,
		Specification: request.Specification,
		Price:         request.Price,
		Unit:          request.Unit,
		Indemnity:     request.Indemnity,
		StockTotal:    request.Stock,
		StockUsed:     0,
		Remark:        request.Remark,
	}

	if t := data.Database.Create(&assets); t.RowsAffected <= 0 {
		response.Fail(ctx, "添加失败")
		return
	}

	response.Success(ctx)
}

func DoDeviceByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID获取失败")
		return
	}

	var request asset.DoDeviceByUpdate
	if err = ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var count int64 = 0
	data.Database.Model(model.DorAssetCategory{}).Where("id=?", request.Category).Where("is_enable=?", constant.IsEnableYes).Count(&count)

	if count <= 0 {
		response.NotFound(ctx, "类型不存在")
		return
	}

	var assets model.DorDevice
	data.Database.Find(&assets, id)
	if assets.Id <= 0 {
		response.NotFound(ctx, "资源不存在")
		return
	}

	if assets.StockUsed > request.Stock {
		response.Fail(ctx, "资源库存不能小于已用量")
		return
	}

	assets.CategoryId = request.Category
	assets.No = request.No
	assets.Name = request.Name
	assets.Specification = request.Specification
	assets.Price = request.Price
	assets.Unit = request.Unit
	assets.Indemnity = request.Indemnity
	assets.StockTotal = request.Stock
	assets.Remark = request.Remark

	if t := data.Database.Save(&assets); t.RowsAffected <= 0 {
		response.Fail(ctx, "修改失败")
		return
	}

	response.Success(ctx)
}

func DoDeviceByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID获取失败")
		return
	}

	var assets model.DorDevice
	data.Database.Find(&assets, id)
	if assets.Id <= 0 {
		response.NotFound(ctx, "资源不存在")
		return
	}

	if t := data.Database.Delete(&assets); t.RowsAffected <= 0 {
		response.Fail(ctx, "删除失败")
		return
	}

	response.Success(ctx)
}

func ToDeviceByPaginate(ctx *gin.Context) {

	var request asset.ToDeviceByPaginate
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	responses := response.Paginate{
		Page: request.GetPage(),
		Size: request.GetSize(),
		Data: make([]any, 0),
	}

	tx := data.Database

	if request.Keyword != "" && request.Type == "no" {
		tx = tx.Where("`no`=?", request.Keyword)
	} else if request.Keyword != "" {
		tx = tx.Where("`name` like ?", "%"+request.Keyword+"%")
	}

	tc := tx

	tc.Model(model.DorDevice{}).Count(&responses.Total)

	if responses.Total > 0 {

		var assets []model.DorDevice
		tx.Preload("Category").Order("`id` desc").Offset(request.GetOffset()).Limit(request.GetLimit()).Find(&assets)

		for _, item := range assets {
			responses.Data = append(responses.Data, asset2.ToDeviceByPaginate{
				Id:            item.Id,
				Category:      item.Category.Name,
				CategoryId:    item.Category.Id,
				Name:          item.Name,
				No:            item.No,
				Specification: item.Specification,
				Price:         item.Price,
				Unit:          item.Unit,
				Indemnity:     item.Indemnity,
				StockTotal:    item.StockTotal,
				StockUsed:     item.StockUsed,
				Remark:        item.Remark,
				CreatedAt:     item.CreatedAt.ToDateTimeString(),
			})
		}
	}

	response.SuccessByPaginate(ctx, responses)
}

func ToDeviceByOnline(ctx *gin.Context) {

	var request asset.ToDeviceByOnline
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	responses := make([]any, 0)

	var devices []model.DorDevice
	data.Database.Where("category_id=?", request.Category).Find(&devices)

	for _, item := range devices {
		responses = append(responses, asset2.ToDeviceByOnline{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	response.SuccessByList(ctx, responses)
}
