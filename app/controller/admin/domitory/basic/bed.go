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

func DoBedByCreate(ctx *gin.Context) {

	var request basic.DoBedByCreate

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var room model.DorRoom

	app.Database.Find(&room, "`id`=? and `is_enable`=?", request.Room, constant.IsEnableYes)

	if room.Id <= 0 {
		response.NotFound(ctx, "房间不存在")
		return
	}

	if room.IsPublic == model.DorBedIsPublicYes {
		response.Fail(ctx, "该房间为公共区域，无法添加")
		return
	}

	bed := model.DorBed{
		BuildingId: room.BuildingId,
		FloorId:    room.FloorId,
		RoomId:     room.Id,
		TypeId:     room.TypeId,
		Name:       request.Name,
		Order:      request.Order,
		IsEnable:   request.IsEnable,
		IsPublic:   request.IsPublic,
	}

	if app.Database.Create(&bed); room.Id <= 0 {
		response.Fail(ctx, "添加失败")
		return
	}

	response.Success(ctx)
}

func DoBedByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil || id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID获取失败")
		return
	}

	var request basic.DoBedByUpdate

	if err = ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var bed model.DorBed

	if app.Database.Find(&bed, id); bed.Id <= 0 {
		response.NotFound(ctx, "未找到该床位")
		return
	}

	if bed.IsEnable != request.IsEnable {
		var peoples int64 = 0
		app.Database.Model(model.DorPeople{}).Where("`bed_id`=? and `status`=?", bed.Id, model.DorPeopleStatusLive).Count(&peoples)
		if peoples > 0 {
			response.Fail(ctx, "该床位已有人入住，无法上下架")
			return
		}
	}

	bed.Name = request.Name
	bed.Order = request.Order
	bed.IsEnable = request.IsEnable

	if t := app.Database.Save(&bed); t.RowsAffected <= 0 {
		response.Fail(ctx, "修改失败")
		return
	}

	response.Success(ctx)
}

func DoBedByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil || id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID获取失败")
		return
	}

	var bed model.DorBed

	if app.Database.Find(&bed, id); bed.Id <= 0 {
		response.NotFound(ctx, "未找到该床位")
		return
	}

	var peoples int64 = 0

	app.Database.Model(model.DorPeople{}).Where("`bed_id`=? and `status`=?", bed.Id, model.DorPeopleStatusLive).Count(&peoples)

	if peoples > 0 {
		response.Fail(ctx, "该床位已有人入住，无法删除")
		return
	}

	if t := app.Database.Delete(&bed); t.RowsAffected <= 0 {
		response.Fail(ctx, "删除失败")
		return
	}

	response.Success(ctx)
}

func DoBedByEnable(ctx *gin.Context) {

	var request basic.DoBedByEnable

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var bed model.DorBed

	if app.Database.Find(&bed, request.Id); bed.Id <= 0 {
		response.NotFound(ctx, "未找到该床位")
		return
	}

	var peoples int64 = 0

	app.Database.Model(model.DorPeople{}).Where("`bed_id`=? and `status`=?", bed.Id, model.DorPeopleStatusLive).Count(&peoples)

	if peoples > 0 {
		response.Fail(ctx, "该床位已有人入住，无法上下架")
		return
	}

	bed.IsEnable = request.IsEnable

	if t := app.Database.Save(&bed); t.RowsAffected <= 0 {
		response.Fail(ctx, "启禁失败")
		return
	}

	response.Success(ctx)
}

func ToBedByPaginate(ctx *gin.Context) {

	var request basic.ToBedByPaginate

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	responses := response.Paginate[res.ToBedByPaginate]{
		Page: request.GetPage(),
		Size: request.GetSize(),
		Data: make([]res.ToBedByPaginate, 0),
	}

	tx := app.Database

	if request.Room > 0 {
		tx = tx.Where("`room_id`", request.Floor)
	} else if request.Floor > 0 {
		tx = tx.Where("`floor_id`", request.Floor)
	} else if request.Building > 0 {
		tx = tx.Where("`building_id`", request.Building)
	}

	if request.IsPublic > 0 {
		tx = tx.Where("`is_public`=?", request.IsPublic)
	}

	if request.Bed != "" {
		tx = tx.Where("`name` like ?", "%"+request.Bed+"%")
	}

	tc := tx

	tc.Model(&model.DorBed{}).Count(&responses.Total)

	if responses.Total > 0 {

		var beds []model.DorBed

		tx.
			Preload("Building").
			Preload("Floor").
			Preload("Room").
			Order("`order` asc, `id` desc").
			Offset(request.GetOffset()).
			Limit(request.GetLimit()).
			Find(&beds)

		for _, item := range beds {
			responses.Data = append(responses.Data, res.ToBedByPaginate{
				Id:        item.Id,
				Name:      item.Name,
				Building:  item.Building.Name,
				Floor:     item.Floor.Name,
				Room:      item.Room.Name,
				Order:     item.Order,
				IsEnable:  item.IsEnable,
				IsPublic:  item.IsPublic,
				CreatedAt: item.CreatedAt.ToDateTimeString(),
			})
		}
	}

	response.SuccessByPaginate(ctx, responses)
}

func ToBedByOnline(ctx *gin.Context) {

	var request basic.ToBedByOnline

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	responses := make([]res.ToBedByOnline, 0)

	tx := app.Database.Where("`room_id`=? and `is_enable`=?", request.Room, constant.IsEnableYes)

	if request.IsPublic > 0 {
		tx = tx.Where("`is_public`=?", request.IsPublic)
	}

	var beds []model.DorBed

	tx.Order("`order` asc, `id` desc").Find(&beds)

	for _, item := range beds {
		items := res.ToBedByOnline{
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
