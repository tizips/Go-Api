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

func DoRoomByCreate(ctx *gin.Context) {

	var former basic.DoRoomByCreateForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var floor model.DorFloor
	data.Database.Where("is_enable", constant.IsEnableYes).First(&floor, former.Floor)
	if floor.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "楼层不存在",
		})
		return
	}

	var count int64
	data.Database.Model(model.DorType{}).Where("id", former.Type).Where("is_enable", constant.IsEnableYes).Count(&count)
	if count <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "房型不存在",
		})
		return
	}

	room := model.DorRoom{
		BuildingId: floor.BuildingId,
		FloorId:    floor.Id,
		TypeId:     former.Type,
		Name:       former.Name,
		Order:      former.Order,
		IsFurnish:  former.IsFurnish,
		IsEnable:   former.IsEnable,
	}

	if data.Database.Create(&room); room.Id <= 0 {
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

func DoRoomByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "房间ID获取失败",
		})
		return
	}

	var former basic.DoRoomByUpdateForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var count int64
	data.Database.Model(model.DorType{}).Where("id", former.Type).Where("is_enable", constant.IsEnableYes).Count(&count)
	if count <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "房型不存在",
		})
		return
	}

	var room model.DorRoom
	data.Database.First(&room, id)
	if room.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该房间",
		})
		return
	}

	room.Name = former.Name
	room.Order = former.Order
	room.TypeId = former.Type
	room.IsFurnish = former.IsFurnish
	room.IsEnable = former.IsEnable

	if t := data.Database.Save(&room); t.RowsAffected <= 0 {
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

func DoRoomByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "房间ID获取失败",
		})
		return
	}

	var room model.DorRoom
	data.Database.First(&room, id)
	if room.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该房间",
		})
		return
	}

	if t := data.Database.Delete(&room); t.RowsAffected <= 0 {
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

func DoRoomByEnable(ctx *gin.Context) {

	var former basic.DoRoomByEnableForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var room model.DorRoom
	data.Database.First(&room, former.Id)
	if room.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该房间",
		})
		return
	}

	room.IsEnable = former.IsEnable

	if t := data.Database.Save(&room); t.RowsAffected <= 0 {
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

func DoRoomByFurnish(ctx *gin.Context) {

	var former basic.DoRoomByFurnishForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var room model.DorRoom
	data.Database.First(&room, former.Id)
	if room.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该房间",
		})
		return
	}

	room.IsFurnish = former.IsFurnish

	if t := data.Database.Save(&room); t.RowsAffected <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "装修失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})

}

func ToRoomByPaginate(ctx *gin.Context) {

	var former basic.ToRoomByPaginateForm
	if err := ctx.ShouldBindQuery(&former); err != nil {
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

	responses.Data.Size = former.GetSize()
	responses.Data.Page = former.GetPage()
	responses.Data.Data = []interface{}{}

	tx := data.Database

	if former.Floor > 0 {
		tx = tx.Where("floor_id", former.Floor)
	} else if former.Building > 0 {
		tx = tx.Where("building_id", former.Building)
	}

	if former.Room != "" {
		tx = tx.Where("name like ?", "%"+former.Room+"%")
	}

	tc := tx

	tc.Model(&model.DorRoom{}).Count(&responses.Data.Total)

	if responses.Data.Total > 0 {

		var rooms []model.DorRoom

		tx.
			Preload("Building").
			Preload("Floor").
			Preload("Type").
			Order("`order` asc").
			Order("`id` desc").
			Find(&rooms)

		for _, item := range rooms {
			responses.Data.Data = append(responses.Data.Data, basicResponse.ToRoomByPaginateResponse{
				Id:        item.Id,
				Name:      item.Name,
				Building:  item.Building.Name,
				Floor:     item.Floor.Name,
				Type:      item.Type.Name,
				TypeId:    item.TypeId,
				Order:     item.Order,
				IsFurnish: item.IsFurnish,
				IsEnable:  item.IsEnable,
				CreatedAt: item.CreatedAt.ToDateTimeString(),
			})
		}
	}

	ctx.JSON(http.StatusOK, responses)
}

func ToRoomByOnline(ctx *gin.Context) {

	var former basic.ToRoomByOnlineForm
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

	var rooms []model.DorRoom
	data.Database.Where("floor_id", former.Floor).Order("`order` asc").Order("`id` desc").Find(&rooms)

	for _, item := range rooms {
		responses.Data = append(responses.Data, basicResponse.ToRoomByOnlineResponse{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	ctx.JSON(http.StatusOK, responses)
}
