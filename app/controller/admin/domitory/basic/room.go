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

	var former basic.DoRoomByCreateFormer
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
	if floor.IsPublic == model.DorFloorIsPublicYes {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "该楼层为公共区域，无法添加",
		})
		return
	}

	room := model.DorRoom{
		BuildingId: floor.BuildingId,
		FloorId:    floor.Id,
		Name:       former.Name,
		Order:      former.Order,
		IsEnable:   former.IsEnable,
		IsPublic:   former.IsPublic,
	}

	var typ model.DorType

	if former.IsPublic != 1 {
		data.Database.Preload("Beds").Where("is_enable=?", constant.IsEnableYes).First(&typ, former.Type)
		if typ.Id <= 0 {
			ctx.JSON(http.StatusOK, response.Response{
				Code:    40400,
				Message: "房型不存在",
			})
			return
		}

		room.TypeId = former.Type
		room.IsFurnish = former.IsFurnish
	}

	tx := data.Database.Begin()

	if t := tx.Create(&room); t.RowsAffected <= 0 {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "添加失败",
		})
		return
	}

	if len(typ.Beds) > 0 {
		var beds []model.DorBed
		for _, item := range typ.Beds {
			beds = append(beds, model.DorBed{
				BuildingId: room.BuildingId,
				FloorId:    room.FloorId,
				RoomId:     room.Id,
				TypeId:     typ.Id,
				BedId:      item.Id,
				Name:       item.Name,
				IsEnable:   constant.IsEnableYes,
				IsPublic:   item.IsPublic,
			})
		}
		if t := tx.Create(&beds); t.RowsAffected <= 0 {
			tx.Rollback()
			ctx.JSON(http.StatusOK, response.Response{
				Code:    40000,
				Message: "添加失败",
			})
			return
		}
	}

	tx.Commit()

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

	var former basic.DoRoomByUpdateFormer
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
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

	if room.IsPublic != 1 {
		var count int64
		data.Database.Model(model.DorType{}).Where("id", former.Type).Where("is_enable", constant.IsEnableYes).Count(&count)
		if count <= 0 {
			ctx.JSON(http.StatusOK, response.Response{
				Code:    40400,
				Message: "房型不存在",
			})
			return
		}
		room.TypeId = former.Type
		room.IsFurnish = former.IsFurnish
	}

	if room.IsEnable != former.IsEnable {
		var peoples int64 = 0
		data.Database.Model(model.DorPeople{}).Where("room_id=?", room.Id).Where("status=?", model.DorPeopleStatusLive).Count(&peoples)
		if peoples > 0 {
			ctx.JSON(http.StatusOK, response.Response{
				Code:    40400,
				Message: "该房间已有人入住，无法上下架",
			})
			return
		}
	}

	room.Name = former.Name
	room.Order = former.Order
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

	var peoples int64 = 0
	data.Database.Model(model.DorPeople{}).Where("room_id=?", room.Id).Where("status=?", model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "该房间已有人入住，无法删除",
		})
		return
	}

	tx := data.Database.Begin()

	if t := tx.Delete(&room); t.RowsAffected <= 0 {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "删除失败",
		})
		return
	}

	if t := tx.Where("room_id=?", room.Id).Delete(&model.DorBed{}); t.Error != nil {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "删除失败",
		})
		return
	}

	tx.Commit()

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})

}

func DoRoomByEnable(ctx *gin.Context) {

	var former basic.DoRoomByEnableFormer
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

	var peoples int64 = 0
	data.Database.Model(model.DorPeople{}).Where("room_id=?", room.Id).Where("status=?", model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "该房间已有人入住，无法上下架",
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

	var former basic.DoRoomByFurnishFormer
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

	var peoples int64 = 0
	data.Database.Model(model.DorPeople{}).Where("room_id=?", room.Id).Where("status=?", model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "该房间已有人入住，无法上下架",
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

	var query basic.ToRoomByPaginateFormer
	if err := ctx.ShouldBindQuery(&query); err != nil {
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

	if query.Floor > 0 {
		tx = tx.Where("floor_id=?", query.Floor)
	} else if query.Building > 0 {
		tx = tx.Where("building_id", query.Building)
	}

	if query.IsPublic > 0 {
		tx = tx.Where("is_public=?", query.IsPublic)
	}

	if query.Room != "" {
		tx = tx.Where("name like ?", "%"+query.Room+"%")
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
			Offset(query.GetOffset()).
			Limit(query.GetLimit()).
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
				IsPublic:  item.IsPublic,
				CreatedAt: item.CreatedAt.ToDateTimeString(),
			})
		}
	}

	ctx.JSON(http.StatusOK, responses)
}

func ToRoomByOnline(ctx *gin.Context) {

	var query basic.ToRoomByOnlineFormer
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

	tx := data.Database.Where("floor_id=?", query.Floor)
	if query.IsPublic > 0 {
		tx = tx.Where("is_public=?", query.IsPublic)
	}

	var rooms []model.DorRoom
	tx.Order("`order` asc").Order("`id` desc").Find(&rooms)

	for _, item := range rooms {
		items := basicResponse.ToRoomByOnlineResponse{
			Id:   item.Id,
			Name: item.Name,
		}
		if query.WithPublic {
			items.IsPublic = item.IsPublic
		}
		responses.Data = append(responses.Data, items)
	}

	ctx.JSON(http.StatusOK, responses)
}
