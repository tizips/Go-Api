package basic

import (
	"github.com/gin-gonic/gin"
	"saas/app/constant"
	"saas/app/model"
	"saas/app/request/admin/dormitory/basic"
	res "saas/app/response/admin/dormitory/basic"
	"saas/kernel/app"
	"saas/kernel/response"
	"strconv"
)

func DoRoomByCreate(ctx *gin.Context) {

	var request basic.DoRoomByCreate

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var floor model.DorFloor

	if app.Database.Find(&floor, "`id`=? and `is_enable`=?", request.Floor, constant.IsEnableYes); floor.Id <= 0 {
		response.NotFound(ctx, "楼层不存在")
		return
	}

	if floor.IsPublic == model.DorFloorIsPublicYes {
		response.Fail(ctx, "该楼层为公共区域，添加失败")
		return
	}

	room := model.DorRoom{
		BuildingId: floor.BuildingId,
		FloorId:    floor.Id,
		Name:       request.Name,
		Order:      request.Order,
		IsEnable:   request.IsEnable,
		IsPublic:   request.IsPublic,
	}

	var typing model.DorType

	if request.IsPublic != 1 {

		app.Database.Preload("Beds").Find(&typing, "`id`=? and `is_enable`=?", request.Type, constant.IsEnableYes)

		if typing.Id <= 0 {
			response.NotFound(ctx, "房型不存在")
			return
		}

		room.TypeId = request.Type
		room.IsFurnish = request.IsFurnish
	}

	tx := app.Database.Begin()

	if t := tx.Create(&room); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "添加失败")
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
			response.Fail(ctx, "添加失败")
			return
		}
	}

	tx.Commit()

	response.Success(ctx)
}

func DoRoomByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil || id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID获取失败")
		return
	}

	var request basic.DoRoomByUpdate

	if err = ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var room model.DorRoom

	if app.Database.Find(&room, id); room.Id <= 0 {
		response.NotFound(ctx, "未找到该房间")
		return
	}

	if room.IsPublic != 1 {

		if room.TypeId != request.Type {

			var count int64

			app.Database.Model(&model.DorType{}).Where("`id`=? and `is_enable`=?", request.Type, constant.IsEnableYes).Count(&count)

			if count <= 0 {
				response.NotFound(ctx, "房型不存在")
				return
			}

			room.TypeId = request.Type
		}

		room.IsFurnish = request.IsFurnish
	}

	if room.IsEnable != request.IsEnable {

		var peoples int64 = 0

		app.Database.Model(model.DorPeople{}).Where("room_id=? and status=?", room.Id, model.DorPeopleStatusLive).Count(&peoples)

		if peoples > 0 {
			response.Fail(ctx, "该房间已有人入住，无法上下架")
			return
		}
	}

	room.Name = request.Name
	room.Order = request.Order
	room.IsEnable = request.IsEnable

	if t := app.Database.Save(&room); t.RowsAffected <= 0 {
		response.Fail(ctx, "修改失败")
		return
	}

	response.Success(ctx)
}

func DoRoomByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil || id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID获取失败")
		return
	}

	var room model.DorRoom

	if app.Database.Find(&room, id); room.Id <= 0 {
		response.NotFound(ctx, "未找到该房间")
		return
	}

	var peoples int64 = 0

	app.Database.Model(&model.DorPeople{}).Where("`room_id`=? and `status`=?", room.Id, model.DorPeopleStatusLive).Count(&peoples)

	if peoples > 0 {
		response.Fail(ctx, "该房间已有人入住，无法删除")
		return
	}

	tx := app.Database.Begin()

	if t := tx.Delete(&room); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "删除失败")
		return
	}

	if t := tx.Where("`room_id`=?", room.Id).Delete(&model.DorBed{}); t.Error != nil {
		tx.Rollback()
		response.Fail(ctx, "删除失败")
		return
	}

	tx.Commit()

	response.Success(ctx)
}

func DoRoomByEnable(ctx *gin.Context) {

	var request basic.DoRoomByEnable

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var room model.DorRoom

	if app.Database.Find(&room, request.Id); room.Id <= 0 {
		response.NotFound(ctx, "未找到该房间")
		return
	}

	var peoples int64 = 0

	app.Database.Model(model.DorPeople{}).Where("room_id=? and status=?", room.Id, model.DorPeopleStatusLive).Count(&peoples)

	if peoples > 0 {
		response.Fail(ctx, "该房间已有人入住，无法上下架")
		return
	}

	room.IsEnable = request.IsEnable

	if t := app.Database.Save(&room); t.RowsAffected <= 0 {
		response.Fail(ctx, "启禁失败")
		return
	}

	response.Success(ctx)
}

func DoRoomByFurnish(ctx *gin.Context) {

	var request basic.DoRoomByFurnish

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var room model.DorRoom

	if app.Database.Find(&room, request.Id); room.Id <= 0 {
		response.NotFound(ctx, "未找到该房间")
		return
	}

	var peoples int64 = 0

	app.Database.Model(model.DorPeople{}).Where("room_id=? and status=?", room.Id, model.DorPeopleStatusLive).Count(&peoples)

	if peoples > 0 {
		response.Fail(ctx, "该房间已有人入住，无法上下架")
		return
	}

	room.IsFurnish = request.IsFurnish

	if t := app.Database.Save(&room); t.RowsAffected <= 0 {
		response.Fail(ctx, "装修失败")
		return
	}

	response.Success(ctx)
}

func ToRoomByPaginate(ctx *gin.Context) {

	var request basic.ToRoomByPaginate

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	responses := response.Paginate[res.ToRoomByPaginate]{
		Page: request.GetPage(),
		Size: request.GetSize(),
		Data: make([]res.ToRoomByPaginate, 0),
	}

	tx := app.Database

	if request.Floor > 0 {
		tx = tx.Where("`floor_id`=?", request.Floor)
	} else if request.Building > 0 {
		tx = tx.Where("`building_id`", request.Building)
	}

	if request.IsPublic > 0 {
		tx = tx.Where("`is_public`=?", request.IsPublic)
	}

	if request.Room != "" {
		tx = tx.Where("`name` like ?", "%"+request.Room+"%")
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
			Offset(request.GetOffset()).
			Limit(request.GetLimit()).
			Find(&rooms)

		for _, item := range rooms {
			responses.Data = append(responses.Data, res.ToRoomByPaginate{
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

	response.SuccessByPaginate(ctx, responses)
}

func ToRoomByOnline(ctx *gin.Context) {

	var request basic.ToRoomByOnline

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	responses := make([]any, 0)

	tx := app.Database.Where("`floor_id`=? and `is_enable`=?", request.Floor, constant.IsEnableYes)

	if request.IsPublic > 0 {
		tx = tx.Where("`is_public`=?", request.IsPublic)
	}

	var rooms []model.DorRoom
	tx.Order("`order` asc, `id` desc").Find(&rooms)

	for _, item := range rooms {
		items := res.ToRoomByOnline{
			Id:   item.Id,
			Name: item.Name,
		}
		if request.WithPublic {
			items.IsPublic = item.IsPublic
		}
		responses = append(responses, items)
	}

	response.SuccessByData(ctx, responses)
}
