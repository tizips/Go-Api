package basic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"saas/app/constant"
	"saas/app/model"
	"saas/app/request/admin/dormitory/basic"
	res "saas/app/response/admin/dormitory/basic"
	"saas/kernel/app"
	"saas/kernel/response"
	"strconv"
)

func DoTypeByCreate(ctx *gin.Context) {

	var request basic.DoTypeByCreate

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	tx := app.MySQL.Begin()

	typ := model.DorType{
		Name:     request.Name,
		Order:    request.Order,
		IsEnable: request.IsEnable,
	}

	if t := tx.Create(&typ); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "添加失败")
		return
	}

	if len(request.Beds) > 0 {

		beds := make([]model.DorTypeBed, len(request.Beds))

		for index, item := range request.Beds {
			beds[index] = model.DorTypeBed{TypeId: typ.Id, Name: item.Name, IsPublic: item.IsPublic}
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

func DoTypeByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil || id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID获取失败")
		return
	}

	var request basic.DoTypeByUpdate

	if err = ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var typing model.DorType

	if app.MySQL.Find(&typing, id); typing.Id <= 0 {
		response.NotFound(ctx, "未找到该房型")
		return
	}

	if typing.IsEnable != request.IsEnable {

		var peoples int64 = 0

		app.MySQL.Model(&model.DorPeople{}).Where("`type_id`=? and `status`=?", typing.Id, model.DorPeopleStatusLive).Count(&peoples)

		if peoples > 0 {
			response.Fail(ctx, "该房型已有人入住，无法上下架")
			return
		}
	}

	typing.Name = request.Name
	typing.Order = request.Order
	typing.IsEnable = request.IsEnable

	if t := app.MySQL.Save(&typing); t.RowsAffected <= 0 {
		response.Fail(ctx, "修改失败")
		return
	}

	response.Success(ctx)
}

func DoTypeByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil || id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID获取失败")
		return
	}

	var typing model.DorType

	if app.MySQL.Find(&typing, id); typing.Id <= 0 {
		response.NotFound(ctx, "未找到该房型")
		return
	}

	var peoples int64 = 0

	app.MySQL.Model(model.DorPeople{}).Where("`type_id`=? and `status`=?", typing.Id, model.DorPeopleStatusLive).Count(&peoples)

	if peoples > 0 {
		response.Fail(ctx, "该房型已有人入住，无法删除")
		return
	}

	tx := app.MySQL.Begin()

	if t := tx.Delete(&typing); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "删除失败")
		return
	}

	if t := tx.Where("type_id=?", typing.Id).Delete(&model.DorTypeBed{}); t.Error != nil {
		tx.Rollback()
		response.Fail(ctx, "删除失败")
		return
	}

	tx.Commit()

	response.Success(ctx)
}

func DoTypeByEnable(ctx *gin.Context) {

	var request basic.DoTypeByEnable

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var typ model.DorType

	if app.MySQL.Find(&typ, request.Id); typ.Id <= 0 {
		response.NotFound(ctx, "未找到该房型")
		return
	}

	var peoples int64 = 0

	app.MySQL.Model(&model.DorPeople{}).Where("`type_id`=? and `status`=?", typ.Id, model.DorPeopleStatusLive).Count(&peoples)

	if peoples > 0 {
		response.Fail(ctx, "该房型已有人入住，无法上下架")
		return
	}

	typ.IsEnable = request.IsEnable

	if t := app.MySQL.Save(&typ); t.RowsAffected <= 0 {
		response.Fail(ctx, "启禁失败")
		return
	}

	response.Success(ctx)
}

func ToTypeByList(ctx *gin.Context) {

	responses := make([]res.ToTypeByList, 0)

	var types []model.DorType

	app.MySQL.Preload("Beds").Order("`order` asc, `id` desc").Find(&types)

	for _, item := range types {
		items := res.ToTypeByList{
			Id:        item.Id,
			Name:      item.Name,
			Order:     item.Order,
			IsEnable:  item.IsEnable,
			CreatedAt: item.CreatedAt.ToDateTimeString(),
		}
		for _, value := range item.Beds {
			items.Beds = append(items.Beds, res.ToTypeByListOfBed{
				Name:     value.Name,
				IsPublic: value.IsPublic,
			})
		}
		responses = append(responses, items)
	}

	response.SuccessByData(ctx, responses)
}

func ToTypeByOnline(ctx *gin.Context) {

	var request basic.ToTypeByOnline

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	responses := make([]res.ToTypeByOnline, 0)

	var types []model.DorType

	tx := app.MySQL.Where("`is_enable`=?", constant.IsEnableYes)

	if request.WithBed || request.MustBed {
		tx = tx.Preload("Beds")
	}
	if request.MustBed {
		tx = tx.Where("exists (?)", app.MySQL.
			Select("1").
			Table(model.TableDorTypeBed).
			Where(fmt.Sprintf("%s.id=%s.type_id", model.TableDorType, model.TableDorTypeBed)).
			Where(fmt.Sprintf("%s.deleted_at is null", model.TableDorTypeBed)),
		)
	}

	tx.Order("`order` asc, `id` desc").Find(&types)

	for _, item := range types {
		items := res.ToTypeByOnline{
			Id:   item.Id,
			Name: item.Name,
		}
		if request.WithBed || request.MustBed {
			for _, value := range item.Beds {
				items.Beds = append(items.Beds, res.ToTypeByOnlineOfBed{
					Id:   value.Id,
					Name: value.Name,
				})
			}
		}
		responses = append(responses, items)
	}

	response.SuccessByData(ctx, responses)
}
