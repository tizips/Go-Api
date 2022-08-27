package basic

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"saas/app/constant"
	"saas/app/model"
	"saas/app/request/admin/dormitory/basic"
	basicResponse "saas/app/response/admin/dormitory/basic"
	"saas/kernel/app"
	"saas/kernel/response"
	"strconv"
)

func DoFloorByCreate(ctx *gin.Context) {

	var request basic.DoFloorByCreate
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var building model.DorBuilding
	app.MySQL.Where("is_enable=?", constant.IsEnableYes).Find(&building, request.Building)
	if building.Id <= 0 {
		response.NotFound(ctx, "楼栋不存在")
		return
	}
	if building.IsPublic == model.DorBuildingIsPublicYes {
		response.Fail(ctx, "该楼栋为公共区域，添加失败")
		return
	}

	floor := model.DorFloor{
		Name:       request.Name,
		BuildingId: request.Building,
		Order:      request.Order,
		IsEnable:   request.IsEnable,
		IsPublic:   request.IsPublic,
	}

	if app.MySQL.Create(&floor); floor.Id <= 0 {
		response.Fail(ctx, "添加失败")
		return
	}

	response.Success(ctx)
}

func DoFloorByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID获取失败")
		return
	}

	var request basic.DoFloorByUpdate
	if err = ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var floor model.DorFloor
	app.MySQL.Find(&floor, id)
	if floor.Id <= 0 {
		response.NotFound(ctx, "未找到该楼层")
		return
	}

	if request.Building != floor.BuildingId {
		var count int64
		app.MySQL.Model(model.DorBuilding{}).Where("`id`=? and `is_enable`=?", request.Building, constant.IsEnableYes).Count(&count)
		if count <= 0 {
			response.NotFound(ctx, "楼栋不存在")
			return
		}

		floor.BuildingId = request.Building
	}

	if floor.IsEnable != request.IsEnable {
		var peoples int64 = 0
		app.MySQL.Model(model.DorPeople{}).Where("`floor_id`=? and `status`=?", floor.Id, model.DorPeopleStatusLive).Count(&peoples)
		if peoples > 0 {
			response.Fail(ctx, "该楼层已有人入住，无法上下架")
			return
		}
	}

	floor.Name = request.Name
	floor.Order = request.Order
	floor.IsEnable = request.IsEnable

	if t := app.MySQL.Save(&floor); t.RowsAffected <= 0 {
		response.Fail(ctx, "修改失败")
		return
	}

	response.Success(ctx)
}

func DoFloorByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID获取失败")
		return
	}

	var floor model.DorFloor
	app.MySQL.Find(&floor, id)
	if floor.Id <= 0 {
		response.NotFound(ctx, "未找到该楼层")
		return
	}

	var peoples int64 = 0
	app.MySQL.Model(model.DorPeople{}).Where("`floor_id`=? and `status`=?", floor.Id, model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "该楼层已有人入住，无法删除",
		})
		return
	}

	tx := app.MySQL.Begin()

	if t := tx.Delete(&floor); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "删除失败")
		return
	}

	if t := tx.Where("`floor_id`=?", floor.Id).Delete(&model.DorRoom{}); t.Error != nil {
		tx.Rollback()
		response.Fail(ctx, "删除失败")
		return
	}

	if t := tx.Where("`floor_id`=?", floor.Id).Delete(&model.DorBed{}); t.Error != nil {
		tx.Rollback()
		response.Fail(ctx, "删除失败")
		return
	}

	tx.Commit()

	response.Success(ctx)
}

func DoFloorByEnable(ctx *gin.Context) {

	var request basic.DoFloorByEnable
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var floor model.DorFloor
	app.MySQL.Find(&floor, request.Id)
	if floor.Id <= 0 {
		response.NotFound(ctx, "未找到该楼层")
		return
	}

	var peoples int64 = 0
	app.MySQL.Model(model.DorPeople{}).Where("`floor_id`=? and `status`=?", floor.Id, model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		response.Fail(ctx, "该楼层已有人入住，无法上下架")
		return
	}

	floor.IsEnable = request.IsEnable

	if t := app.MySQL.Save(&floor); t.RowsAffected <= 0 {
		response.Fail(ctx, "启禁失败")
		return
	}

	response.Success(ctx)
}

func ToFloorByList(ctx *gin.Context) {

	var request basic.ToFloorByList
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	responses := make([]any, 0)

	var floors []model.DorFloor
	app.MySQL.
		Preload("Building").
		Where("`building_id`=?", request.Building).
		Order("`order` asc, `id` desc").
		Find(&floors)

	for _, item := range floors {
		responses = append(responses, basicResponse.ToFloorByList{
			Id:        item.Id,
			Name:      item.Name,
			Building:  item.Building.Name,
			Order:     item.Order,
			IsEnable:  item.IsEnable,
			IsPublic:  item.IsPublic,
			CreatedAt: item.CreatedAt.ToDateTimeString(),
		})
	}

	response.SuccessByList(ctx, responses)
}

func ToFloorByOnline(ctx *gin.Context) {

	var request basic.ToFloorByOnline
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	responses := make([]any, 0)

	tx := app.MySQL.Where("`building_id`=? and `is_enable`=?", request.Building, constant.IsEnableYes)

	if request.IsPublic > 0 {
		tx = tx.Where("`is_public`=?", request.IsPublic)
	}

	var floors []model.DorFloor
	tx.Order("`order` asc, `id` desc").Order("`id` desc").Find(&floors)

	for _, item := range floors {
		items := basicResponse.ToFloorByOnline{
			Id:       item.Id,
			Name:     item.Name,
			IsPublic: item.IsPublic,
		}
		if request.WithPublic {
			items.IsPublic = item.IsPublic
		}
		responses = append(responses, items)
	}

	response.SuccessByList(ctx, responses)
}
