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
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var building model.DorBuilding
	data.Database.Where("is_enable=?", constant.IsEnableYes).First(&building, former.Building)
	if building.Id <= 0 {
		response.ToResponseByNotFound(ctx, "楼栋不存在")
		return
	}
	if building.IsPublic == model.DorBuildingIsPublicYes {
		response.ToResponseByFail(ctx, "该楼栋为公共区域，添加失败")
		return
	}

	floor := model.DorFloor{
		Name:       former.Name,
		BuildingId: former.Building,
		Order:      former.Order,
		IsEnable:   former.IsEnable,
		IsPublic:   former.IsPublic,
	}

	if data.Database.Create(&floor); floor.Id <= 0 {
		response.ToResponseByFail(ctx, "添加失败")
		return
	}

	response.ToResponseBySuccess(ctx)
}

func DoFloorByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.ToResponseByFailRequestMessage(ctx, "ID获取失败")
		return
	}

	var former basic.DoFloorByUpdateForm
	if err = ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var floor model.DorFloor
	data.Database.First(&floor, id)
	if floor.Id <= 0 {
		response.ToResponseByNotFound(ctx, "未找到该楼层")
		return
	}

	if former.Building != floor.BuildingId {
		var count int64
		data.Database.Model(model.DorBuilding{}).Where("`id`=? and `is_enable`=?", former.Building, constant.IsEnableYes).Count(&count)
		if count <= 0 {
			response.ToResponseByNotFound(ctx, "楼栋不存在")
			return
		}

		floor.BuildingId = former.Building
	}

	if floor.IsEnable != former.IsEnable {
		var peoples int64 = 0
		data.Database.Model(model.DorPeople{}).Where("`floor_id`=? and `status`=?", floor.Id, model.DorPeopleStatusLive).Count(&peoples)
		if peoples > 0 {
			response.ToResponseByFail(ctx, "该楼层已有人入住，无法上下架")
			return
		}
	}

	floor.Name = former.Name
	floor.Order = former.Order
	floor.IsEnable = former.IsEnable

	if t := data.Database.Save(&floor); t.RowsAffected <= 0 {
		response.ToResponseByFail(ctx, "修改失败")
		return
	}

	response.ToResponseBySuccess(ctx)
}

func DoFloorByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.ToResponseByFailRequestMessage(ctx, "ID获取失败")
		return
	}

	var floor model.DorFloor
	data.Database.First(&floor, id)
	if floor.Id <= 0 {
		response.ToResponseByNotFound(ctx, "未找到该楼层")
		return
	}

	var peoples int64 = 0
	data.Database.Model(model.DorPeople{}).Where("`floor_id`=? and `status`=?", floor.Id, model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "该楼层已有人入住，无法删除",
		})
		return
	}

	tx := data.Database.Begin()

	if t := tx.Delete(&floor); t.RowsAffected <= 0 {
		tx.Rollback()
		response.ToResponseByFail(ctx, "删除失败")
		return
	}

	if t := tx.Where("`floor_id`=?", floor.Id).Delete(&model.DorRoom{}); t.Error != nil {
		tx.Rollback()
		response.ToResponseByFail(ctx, "删除失败")
		return
	}

	if t := tx.Where("`floor_id`=?", floor.Id).Delete(&model.DorBed{}); t.Error != nil {
		tx.Rollback()
		response.ToResponseByFail(ctx, "删除失败")
		return
	}

	tx.Commit()

	response.ToResponseBySuccess(ctx)
}

func DoFloorByEnable(ctx *gin.Context) {

	var former basic.DoFloorByEnableForm
	if err := ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var floor model.DorFloor
	data.Database.First(&floor, former.Id)
	if floor.Id <= 0 {
		response.ToResponseByNotFound(ctx, "未找到该楼层")
		return
	}

	var peoples int64 = 0
	data.Database.Model(model.DorPeople{}).Where("`floor_id`=? and `status`=?", floor.Id, model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		response.ToResponseByFail(ctx, "该楼层已有人入住，无法上下架")
		return
	}

	floor.IsEnable = former.IsEnable

	if t := data.Database.Save(&floor); t.RowsAffected <= 0 {
		response.ToResponseByFail(ctx, "启禁失败")
		return
	}

	response.ToResponseBySuccess(ctx)
}

func ToFloorByList(ctx *gin.Context) {

	var former basic.ToFloorByListForm
	if err := ctx.ShouldBindQuery(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	responses := make([]any, 0)

	var floors []model.DorFloor
	data.Database.
		Preload("Building").
		Where("`building_id`=?", former.Building).
		Order("`order` asc, `id` desc").
		Find(&floors)

	for _, item := range floors {
		responses = append(responses, basicResponse.ToFloorByListResponse{
			Id:        item.Id,
			Name:      item.Name,
			Building:  item.Building.Name,
			Order:     item.Order,
			IsEnable:  item.IsEnable,
			IsPublic:  item.IsPublic,
			CreatedAt: item.CreatedAt.ToDateTimeString(),
		})
	}

	response.ToResponseBySuccessList(ctx, responses)
}

func ToFloorByOnline(ctx *gin.Context) {

	var query basic.ToFloorByOnlineForm
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	responses := make([]any, 0)

	tx := data.Database.Where("`building_id`=? and `is_enable`=?", query.Building, constant.IsEnableYes)

	if query.IsPublic > 0 {
		tx = tx.Where("`is_public`=?", query.IsPublic)
	}

	var floors []model.DorFloor
	tx.Order("`order` asc, `id` desc").Order("`id` desc").Find(&floors)

	for _, item := range floors {
		items := basicResponse.ToFloorByOnlineResponse{
			Id:       item.Id,
			Name:     item.Name,
			IsPublic: item.IsPublic,
		}
		if query.WithPublic {
			items.IsPublic = item.IsPublic
		}
		responses = append(responses, items)
	}

	response.ToResponseBySuccessList(ctx, responses)
}
