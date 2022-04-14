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

func DoBuildingByCreate(ctx *gin.Context) {

	var former basic.DoBuildingByCreateFormer
	if err := ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	building := model.DorBuilding{
		Name:     former.Name,
		Order:    former.Order,
		IsEnable: former.IsEnable,
		IsPublic: former.IsPublic,
	}

	if data.Database.Create(&building); building.Id <= 0 {
		response.ToResponseByFail(ctx, "添加失败")
		return
	}

	response.ToResponseBySuccess(ctx)
}

func DoBuildingByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.ToResponseByFailRequestMessage(ctx, "ID获取失败")
		return
	}

	var former basic.DoBuildingByUpdateFormer
	if err = ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var building model.DorBuilding
	data.Database.First(&building, id)
	if building.Id <= 0 {
		response.ToResponseByNotFound(ctx, "未找到该楼栋")
		return
	}

	if building.IsEnable != former.IsEnable {
		var peoples int64 = 0
		data.Database.Model(model.DorPeople{}).Where("`building_id`=? and `status`=?", building.Id, model.DorPeopleStatusLive).Count(&peoples)
		if peoples > 0 {
			response.ToResponseByFail(ctx, "该楼栋已有人入住，无法上下架")
			return
		}
	}

	building.Name = former.Name
	building.Order = former.Order
	building.IsEnable = former.IsEnable

	if t := data.Database.Save(&building); t.RowsAffected <= 0 {
		response.ToResponseByFail(ctx, "修改失败")
		return
	}

	response.ToResponseBySuccess(ctx)
}

func DoBuildingByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.ToResponseByFailRequestMessage(ctx, "ID获取失败")
		return
	}

	var building model.DorBuilding
	data.Database.First(&building, id)
	if building.Id <= 0 {
		response.ToResponseByNotFound(ctx, "未找到该楼栋")
		return
	}

	var peoples int64 = 0
	data.Database.Model(model.DorPeople{}).Where("`building_id`=? and `status`=?", building.Id, model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		response.ToResponseByFail(ctx, "该楼栋已有人入住，无法删除")
		return
	}

	tx := data.Database.Begin()

	if t := tx.Delete(&building); t.RowsAffected <= 0 {
		tx.Rollback()
		response.ToResponseByFail(ctx, "删除失败")
		return
	}

	if t := tx.Where("building_id=?", building.Id).Delete(&model.DorRoom{}); t.Error != nil {
		tx.Rollback()
		response.ToResponseByFail(ctx, "删除失败")
		return
	}

	if t := tx.Where("building_id=?", building.Id).Delete(&model.DorRoom{}); t.Error != nil {
		tx.Rollback()
		response.ToResponseByFail(ctx, "删除失败")
		return
	}

	if t := tx.Where("building_id=?", building.Id).Delete(&model.DorBed{}); t.Error != nil {
		tx.Rollback()
		response.ToResponseByFail(ctx, "删除失败")
		return
	}

	tx.Commit()

	response.ToResponseBySuccess(ctx)
}

func DoBuildingByEnable(ctx *gin.Context) {

	var former basic.DoBuildingByEnableFormer
	if err := ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var building model.DorBuilding
	data.Database.First(&building, former.Id)
	if building.Id <= 0 {
		response.ToResponseByNotFound(ctx, "未找到该楼栋")
		return
	}

	var peoples int64 = 0
	data.Database.Model(model.DorPeople{}).Where("`building_id`=? and `status`=?", building.Id, model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		response.ToResponseByFail(ctx, "该楼栋已有人入住，无法上下架")
		return
	}

	building.IsEnable = former.IsEnable

	if t := data.Database.Save(&building); t.RowsAffected <= 0 {
		response.ToResponseByFail(ctx, "启禁失败")
		return
	}

	response.ToResponseBySuccess(ctx)
}

func ToBuildingByList(ctx *gin.Context) {

	responses := make([]any, 0)

	var buildings []model.DorBuilding
	data.Database.Order("`order` asc, `id` desc").Find(&buildings)

	for _, item := range buildings {
		responses = append(responses, basicResponse.ToBuildingByListResponse{
			Id:        item.Id,
			Name:      item.Name,
			Order:     item.Order,
			IsEnable:  item.IsEnable,
			IsPublic:  item.IsPublic,
			CreatedAt: item.CreatedAt.ToDateTimeString(),
		})
	}

	response.ToResponseBySuccessList(ctx, responses)
}

func ToBuildingByOnline(ctx *gin.Context) {

	var query basic.ToBuildingByOnlineFormer
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	responses := make([]any, 0)

	tx := data.Database.Where("`is_enable`=?", constant.IsEnableYes)

	if query.IsPublic > 0 {
		tx = tx.Where("`is_public`=?", query.IsPublic)
	}

	var buildings []model.DorBuilding
	tx.Order("`order` asc, `id` desc").Find(&buildings)

	for _, item := range buildings {
		items := basicResponse.ToBuildingByOnlineResponse{
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
