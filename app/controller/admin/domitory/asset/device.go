package asset

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"saas/app/constant"
	"saas/app/form/admin/dormitory/asset"
	"saas/app/model"
	asset2 "saas/app/response/admin/dormitory/asset"
	"saas/kernel/data"
	"saas/kernel/response"
	"strconv"
)

func DoDeviceByCreate(ctx *gin.Context) {

	var former asset.DoDeviceByCreateForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var count int64 = 0
	data.Database.Model(model.DorAssetCategory{}).Where("id=?", former.Category).Where("is_enable=?", constant.IsEnableYes).Count(&count)

	if count <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "类型不存在",
		})
		return
	}

	assets := model.DorDevice{
		CategoryId:    former.Category,
		No:            former.No,
		Name:          former.Name,
		Specification: former.Specification,
		Price:         former.Price,
		Unit:          former.Unit,
		Indemnity:     former.Indemnity,
		StockTotal:    former.Stock,
		StockUsed:     0,
		Remark:        former.Remark,
	}

	if t := data.Database.Create(&assets); t.RowsAffected <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "添加失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})
}

func DoDeviceByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "资源ID获取失败",
		})
		return
	}

	var former asset.DoDeviceByUpdateForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    20000,
			Message: err.Error(),
		})
		return
	}

	var count int64 = 0
	data.Database.Model(model.DorAssetCategory{}).Where("id=?", former.Category).Where("is_enable=?", constant.IsEnableYes).Count(&count)

	if count <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "类型不存在",
		})
		return
	}

	var assets model.DorDevice
	data.Database.Find(&assets, id)
	if assets.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "资源不存在",
		})
		return
	}

	if assets.StockUsed > former.Stock {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "资源库存不能小于已用量",
		})
		return
	}

	assets.CategoryId = former.Category
	assets.No = former.No
	assets.Name = former.Name
	assets.Specification = former.Specification
	assets.Price = former.Price
	assets.Unit = former.Unit
	assets.Indemnity = former.Indemnity
	assets.StockTotal = former.Stock
	assets.Remark = former.Remark

	if t := data.Database.Save(&assets); t.RowsAffected <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "修改失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})
}

func DoDeviceByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    20000,
			Message: "资源ID获取失败",
		})
		return
	}

	var assets model.DorDevice
	data.Database.Find(&assets, id)
	if assets.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "资源不存在",
		})
		return
	}

	if t := data.Database.Delete(&assets); t.RowsAffected <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "删除失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})
}

func ToDeviceByPaginate(ctx *gin.Context) {

	var query asset.ToDeviceByPaginateForm
	if err := ctx.ShouldBind(&query); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	responses := response.Paginate{
		Code:    20000,
		Message: "Success",
	}

	responses.Data.Size = query.GetSize()
	responses.Data.Page = query.GetPage()
	responses.Data.Data = []any{}

	tx := data.Database

	if query.Keyword != "" && query.Type == "no" {
		tx = tx.Where("no=?", query.Keyword)
	} else if query.Keyword != "" {
		tx = tx.Where("name like ?", "%"+query.Keyword+"%")
	}

	tc := tx

	tc.Model(model.DorDevice{}).Count(&responses.Data.Total)

	if responses.Data.Total > 0 {

		var assets []model.DorDevice
		tx.Preload("Category").Order("`id` desc").Offset(query.GetOffset()).Limit(query.GetLimit()).Find(&assets)

		for _, item := range assets {
			responses.Data.Data = append(responses.Data.Data, asset2.ToDeviceByPaginateResponse{
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

	ctx.JSON(http.StatusOK, responses)

}

func ToDeviceByOnline(ctx *gin.Context) {

	var query asset.ToDeviceByOnlineForm
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	responses := response.Responses{
		Code:    20000,
		Message: "Success",
		Data:    []any{},
	}

	var devices []model.DorDevice
	data.Database.Where("category_id=?", query.Category).Find(&devices)

	for _, item := range devices {
		responses.Data = append(responses.Data, asset2.ToDeviceByOnlineResponse{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	ctx.JSON(http.StatusOK, responses)
}
