package basic

import (
	"github.com/gin-gonic/gin"
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
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var floor model.DorFloor
	data.Database.Where("is_enable", constant.IsEnableYes).First(&floor, former.Floor)
	if floor.Id <= 0 {
		response.ToResponseByNotFound(ctx, "楼层不存在")
		return
	}
	if floor.IsPublic == model.DorFloorIsPublicYes {
		response.ToResponseByFail(ctx, "该楼层为公共区域，添加失败")
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

	var typing model.DorType

	if former.IsPublic != 1 {
		data.Database.Preload("Beds").Where("`is_enable`=?", constant.IsEnableYes).First(&typing, former.Type)
		if typing.Id <= 0 {
			response.ToResponseByNotFound(ctx, "房型不存在")
			return
		}

		room.TypeId = former.Type
		room.IsFurnish = former.IsFurnish
	}

	tx := data.Database.Begin()

	if t := tx.Create(&room); t.RowsAffected <= 0 {
		tx.Rollback()
		response.ToResponseByFail(ctx, "添加失败")
		return
	}

	if len(typing.Beds) > 0 {
		var beds []model.DorBed
		for _, item := range typing.Beds {
			beds = append(beds, model.DorBed{
				BuildingId: room.BuildingId,
				FloorId:    room.FloorId,
				RoomId:     room.Id,
				TypeId:     typing.Id,
				BedId:      item.Id,
				Name:       item.Name,
				IsEnable:   constant.IsEnableYes,
				IsPublic:   item.IsPublic,
			})
		}
		if t := tx.Create(&beds); t.RowsAffected <= 0 {
			tx.Rollback()
			response.ToResponseByFail(ctx, "添加失败")
			return
		}
	}

	tx.Commit()

	response.ToResponseBySuccess(ctx)
}

func DoRoomByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.ToResponseByFailRequestMessage(ctx, "ID获取失败")
		return
	}

	var former basic.DoRoomByUpdateFormer
	if err = ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var room model.DorRoom
	data.Database.First(&room, id)
	if room.Id <= 0 {
		response.ToResponseByNotFound(ctx, "未找到该房间")
		return
	}

	if room.IsPublic != 1 {
		if room.TypeId != former.Type {
			var count int64
			data.Database.Model(&model.DorType{}).Where("`id`=? and `is_enable`=?", former.Type, constant.IsEnableYes).Count(&count)
			if count <= 0 {
				response.ToResponseByNotFound(ctx, "房型不存在")
				return
			}
			room.TypeId = former.Type
		}
		room.IsFurnish = former.IsFurnish
	}

	if room.IsEnable != former.IsEnable {
		var peoples int64 = 0
		data.Database.Model(model.DorPeople{}).Where("room_id=? and status=?", room.Id, model.DorPeopleStatusLive).Count(&peoples)
		if peoples > 0 {
			response.ToResponseByFail(ctx, "该房间已有人入住，无法上下架")
			return
		}
	}

	room.Name = former.Name
	room.Order = former.Order
	room.IsEnable = former.IsEnable

	if t := data.Database.Save(&room); t.RowsAffected <= 0 {
		response.ToResponseByFail(ctx, "修改失败")
		return
	}

	response.ToResponseBySuccess(ctx)
}

func DoRoomByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.ToResponseByFailRequestMessage(ctx, "ID获取失败")
		return
	}

	var room model.DorRoom
	data.Database.First(&room, id)
	if room.Id <= 0 {
		response.ToResponseByNotFound(ctx, "未找到该房间")
		return
	}

	var peoples int64 = 0
	data.Database.Model(&model.DorPeople{}).Where("`room_id`=? and `status`=?", room.Id, model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		response.ToResponseByFail(ctx, "该房间已有人入住，无法删除")
		return
	}

	tx := data.Database.Begin()

	if t := tx.Delete(&room); t.RowsAffected <= 0 {
		tx.Rollback()
		response.ToResponseByFail(ctx, "删除失败")
		return
	}

	if t := tx.Where("`room_id`=?", room.Id).Delete(&model.DorBed{}); t.Error != nil {
		tx.Rollback()
		response.ToResponseByFail(ctx, "删除失败")
		return
	}

	tx.Commit()

	response.ToResponseBySuccess(ctx)
}

func DoRoomByEnable(ctx *gin.Context) {

	var former basic.DoRoomByEnableFormer
	if err := ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var room model.DorRoom
	data.Database.First(&room, former.Id)
	if room.Id <= 0 {
		response.ToResponseByNotFound(ctx, "未找到该房间")
		return
	}

	var peoples int64 = 0
	data.Database.Model(model.DorPeople{}).Where("room_id=? and status=?", room.Id, model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		response.ToResponseByFail(ctx, "该房间已有人入住，无法上下架")
		return
	}

	room.IsEnable = former.IsEnable

	if t := data.Database.Save(&room); t.RowsAffected <= 0 {
		response.ToResponseByFail(ctx, "启禁失败")
		return
	}

	response.ToResponseBySuccess(ctx)
}

func DoRoomByFurnish(ctx *gin.Context) {

	var former basic.DoRoomByFurnishFormer
	if err := ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var room model.DorRoom
	data.Database.First(&room, former.Id)
	if room.Id <= 0 {
		response.ToResponseByNotFound(ctx, "未找到该房间")
		return
	}

	var peoples int64 = 0
	data.Database.Model(model.DorPeople{}).Where("room_id=? and status=?", room.Id, model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		response.ToResponseByFail(ctx, "该房间已有人入住，无法上下架")
		return
	}

	room.IsFurnish = former.IsFurnish

	if t := data.Database.Save(&room); t.RowsAffected <= 0 {
		response.ToResponseByFail(ctx, "装修失败")
		return
	}

	response.ToResponseBySuccess(ctx)
}

func ToRoomByPaginate(ctx *gin.Context) {

	var query basic.ToRoomByPaginateFormer
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	responses := response.Paginate{
		Page: query.GetPage(),
		Size: query.GetSize(),
		Data: make([]any, 0),
	}

	tx := data.Database

	if query.Floor > 0 {
		tx = tx.Where("`floor_id`=?", query.Floor)
	} else if query.Building > 0 {
		tx = tx.Where("`building_id`", query.Building)
	}

	if query.IsPublic > 0 {
		tx = tx.Where("`is_public`=?", query.IsPublic)
	}

	if query.Room != "" {
		tx = tx.Where("`name` like ?", "%"+query.Room+"%")
	}

	tc := tx

	tc.Model(&model.DorRoom{}).Count(&responses.Total)

	if responses.Total > 0 {

		var rooms []model.DorRoom

		tx.
			Preload("Building").
			Preload("Floor").
			Preload("Type").
			Order("`order` asc, `id` desc").
			Offset(query.GetOffset()).
			Limit(query.GetLimit()).
			Find(&rooms)

		for _, item := range rooms {
			responses.Data = append(responses.Data, basicResponse.ToRoomByPaginateResponse{
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

	response.ToResponseBySuccessPaginate(ctx, responses)
}

func ToRoomByOnline(ctx *gin.Context) {

	var query basic.ToRoomByOnlineFormer
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	responses := make([]any, 0)

	tx := data.Database.Where("`floor_id`=? and `is_enable`=?", query.Floor, constant.IsEnableYes)

	if query.IsPublic > 0 {
		tx = tx.Where("`is_public`=?", query.IsPublic)
	}

	var rooms []model.DorRoom
	tx.Order("`order` asc, `id` desc").Find(&rooms)

	for _, item := range rooms {
		items := basicResponse.ToRoomByOnlineResponse{
			Id:   item.Id,
			Name: item.Name,
		}
		if query.WithPublic {
			items.IsPublic = item.IsPublic
		}
		responses = append(responses, items)
	}

	response.ToResponseBySuccessList(ctx, responses)
}
