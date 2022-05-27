package basic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"saas/app/constant"
	"saas/app/model"
	"saas/app/request/admin/dormitory/basic"
	basicResponse "saas/app/response/admin/dormitory/basic"
	"saas/kernel/data"
	"saas/kernel/response"
	"strconv"
)

func DoTypeByCreate(ctx *gin.Context) {

	var request basic.DoTypeByCreate
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	tx := data.Database.Begin()

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
	data.Database.Find(&typing, id)
	if typing.Id <= 0 {
		response.NotFound(ctx, "未找到该房型")
		return
	}

	if typing.IsEnable != request.IsEnable {
		var peoples int64 = 0
		data.Database.Model(&model.DorPeople{}).Where("`type_id`=? and `status`=?", typing.Id, model.DorPeopleStatusLive).Count(&peoples)
		if peoples > 0 {
			response.Fail(ctx, "该房型已有人入住，无法上下架")
			return
		}
	}

	typing.Name = request.Name
	typing.Order = request.Order
	typing.IsEnable = request.IsEnable

	if t := data.Database.Save(&typing); t.RowsAffected <= 0 {
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
	data.Database.Find(&typing, id)
	if typing.Id <= 0 {
		response.NotFound(ctx, "未找到该房型")
		return
	}

	var peoples int64 = 0
	data.Database.Model(model.DorPeople{}).Where("`type_id`=? and `status`=?", typing.Id, model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		response.Fail(ctx, "该房型已有人入住，无法删除")
		return
	}

	tx := data.Database.Begin()

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
	data.Database.Find(&typ, request.Id)
	if typ.Id <= 0 {
		response.NotFound(ctx, "未找到该房型")
		return
	}

	var peoples int64 = 0
	data.Database.Model(&model.DorPeople{}).Where("`type_id`=? and `status`=?", typ.Id, model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		response.Fail(ctx, "该房型已有人入住，无法上下架")
		return
	}

	typ.IsEnable = request.IsEnable

	if t := data.Database.Save(&typ); t.RowsAffected <= 0 {
		response.Fail(ctx, "启禁失败")
		return
	}

	response.Success(ctx)
}

func ToTypeByList(ctx *gin.Context) {

	responses := make([]any, 0)

	var types []model.DorType
	data.Database.Preload("Beds").Order("`order` asc, `id` desc").Find(&types)

	for _, item := range types {
		items := basicResponse.ToTypeByList{
			Id:        item.Id,
			Name:      item.Name,
			Order:     item.Order,
			IsEnable:  item.IsEnable,
			CreatedAt: item.CreatedAt.ToDateTimeString(),
		}
		for _, value := range item.Beds {
			items.Beds = append(items.Beds, basicResponse.ToTypeByListOfBed{
				Name:     value.Name,
				IsPublic: value.IsPublic,
			})
		}
		responses = append(responses, items)
	}

	response.SuccessByList(ctx, responses)
}

func ToTypeByOnline(ctx *gin.Context) {

	var request basic.ToTypeByOnline
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	responses := make([]any, 0)

	var types []model.DorType

	tx := data.Database.Where("`is_enable`=?", constant.IsEnableYes)

	if request.WithBed || request.MustBed {
		tx = tx.Preload("Beds")
	}
	if request.MustBed {
		tx = tx.Where("exists (?)", data.Database.
			Select("1").
			Table(model.TableDorTypeBed).
			Where(fmt.Sprintf("%s.id=%s.type_id", model.TableDorType, model.TableDorTypeBed)).
			Where(fmt.Sprintf("%s.deleted_at is null", model.TableDorTypeBed)),
		)
	}

	tx.Order("`order` asc, `id` desc").Find(&types)

	for _, item := range types {
		items := basicResponse.ToTypeByOnline{
			Id:   item.Id,
			Name: item.Name,
		}
		if request.WithBed || request.MustBed {
			for _, value := range item.Beds {
				items.Beds = append(items.Beds, basicResponse.ToTypeByOnlineOfBed{
					Id:   value.Id,
					Name: value.Name,
				})
			}
		}
		responses = append(responses, items)
	}

	response.SuccessByList(ctx, responses)
}
