package asset

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	category := model.DorAssetCategory{
		Name:     former.Name,
		Order:    former.Order,
		IsEnable: former.IsEnable,
	}

	if data.Database.Create(&category); category.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "添加失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})

}

func DoCategoryByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "类型ID获取失败",
		})
		return
	}

	var former asset.DoCategoryByUpdateForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var category model.DorAssetCategory
	data.Database.First(&category, id)
	if category.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该类型",
		})
		return
	}

	category.Name = former.Name
	category.Order = former.Order
	category.IsEnable = former.IsEnable

	if t := data.Database.Save(&category); t.RowsAffected <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "修改失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})

}

func DoCategoryByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "类型ID获取失败",
		})
		return
	}

	var category model.DorAssetCategory
	data.Database.First(&category, id)
	if category.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该类型",
		})
		return
	}

	if t := data.Database.Delete(&category); t.RowsAffected <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "删除失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})

}

func DoCategoryByEnable(ctx *gin.Context) {

	var former asset.DoCategoryByEnableForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var category model.DorAssetCategory
	data.Database.First(&category, former.Id)
	if category.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该类型",
		})
		return
	}

	category.IsEnable = former.IsEnable

	if t := data.Database.Save(&category); t.RowsAffected <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "启禁失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})

}

func ToCategoryByList(ctx *gin.Context) {

	responses := response.Responses{
		Code:    20000,
		Message: "Success",
		Data:    []any{},
	}

	var categories []model.DorAssetCategory
	data.Database.Order("`order` asc").Order("`id` desc").Find(&categories)

	for _, item := range categories {
		responses.Data = append(responses.Data, assetResponse.ToCategoryByListResponse{
			Id:        item.Id,
			Name:      item.Name,
			Order:     item.Order,
			IsEnable:  item.IsEnable,
			CreatedAt: item.CreatedAt.ToDateTimeString(),
		})
	}

	ctx.JSON(http.StatusOK, responses)
}

func ToCategoryByOnline(ctx *gin.Context) {

	responses := response.Responses{
		Code:    20000,
		Message: "Success",
		Data:    []any{},
	}

	var categories []model.DorAssetCategory
	data.Database.Where("is_enable=?", constant.IsEnableYes).Order("`order` asc").Order("`id` desc").Find(&categories)

	for _, item := range categories {
		responses.Data = append(responses.Data, assetResponse.ToCategoryByOnlineResponse{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	ctx.JSON(http.StatusOK, responses)
}
