package basic

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"saas/app/constant"
	"saas/app/form/admin/dormitory/basic"
	"saas/app/model"
	basicResponse "saas/app/response/admin/dormitory/basic"
	"saas/kernel/data"
	"saas/kernel/response"
	"strconv"
)

func DoFloorByCreate(ctx *gin.Context) {

	var former basic.DoFloorByCreateForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var count int64
	data.Database.Model(model.DorBuilding{}).Where("id", former.Building).Where("is_enable", constant.IsEnableYes).Count(&count)
	if count <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "楼栋不存在",
		})
		return
	}

	floor := model.DorFloor{
		Name:       former.Name,
		BuildingId: former.Building,
		Order:      former.Order,
		IsEnable:   former.IsEnable,
	}

	if data.Database.Create(&floor); floor.Id <= 0 {
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

func DoFloorByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "楼层ID获取失败",
		})
		return
	}

	var former basic.DoFloorByUpdateForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var count int64
	data.Database.Model(model.DorBuilding{}).Where("id", former.Building).Where("is_enable", constant.IsEnableYes).Count(&count)
	if count <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "楼栋不存在",
		})
		return
	}

	var floor model.DorFloor
	data.Database.First(&floor, id)
	if floor.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该楼栋",
		})
		return
	}

	floor.BuildingId = former.Building
	floor.Name = former.Name
	floor.Order = former.Order
	floor.IsEnable = former.IsEnable

	if t := data.Database.Save(&floor); t.RowsAffected <= 0 {
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

func DoFloorByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "楼栋ID获取失败",
		})
		return
	}

	var floor model.DorFloor
	data.Database.First(&floor, id)
	if floor.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该楼层",
		})
		return
	}

	if t := data.Database.Delete(&floor); t.RowsAffected <= 0 {
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

func DoFloorByEnable(ctx *gin.Context) {

	var former basic.DoFloorByEnableForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var floor model.DorFloor
	data.Database.First(&floor, former.Id)
	if floor.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该楼层",
		})
		return
	}

	floor.IsEnable = former.IsEnable

	if t := data.Database.Save(&floor); t.RowsAffected <= 0 {
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

func ToFloorByList(ctx *gin.Context) {

	var former basic.ToFloorByListForm
	if err := ctx.ShouldBindQuery(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	responses := response.Responses{
		Code:    20000,
		Message: "Success",
		Data:    []interface{}{},
	}

	var floors []model.DorFloor
	data.Database.
		Preload("Building").
		Where("building_id", former.Building).
		Order("`order` asc").
		Order("`id` desc").
		Find(&floors)

	for _, item := range floors {
		responses.Data = append(responses.Data, basicResponse.ToFloorByListResponse{
			Id:        item.Id,
			Name:      item.Name,
			Building:  item.Building.Name,
			Order:     item.Order,
			IsEnable:  item.IsEnable,
			CreatedAt: item.CreatedAt.ToDateTimeString(),
		})
	}

	ctx.JSON(http.StatusOK, responses)
}

func ToFloorByOnline(ctx *gin.Context) {

	var former basic.ToFloorByOnlineForm
	if err := ctx.ShouldBindQuery(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	responses := response.Responses{
		Code:    20000,
		Message: "Success",
		Data:    []interface{}{},
	}

	var floors []model.DorFloor
	data.Database.Where("building_id", former.Building).Order("`order` asc").Order("`id` desc").Find(&floors)

	for _, item := range floors {
		responses.Data = append(responses.Data, basicResponse.ToFloorByOnlineResponse{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	ctx.JSON(http.StatusOK, responses)
}
