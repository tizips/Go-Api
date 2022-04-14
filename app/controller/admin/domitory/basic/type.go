package basic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"saas/app/constant"
	"saas/app/form/admin/dormitory/basic"
	"saas/app/model"
	basicResponse "saas/app/response/admin/dormitory/basic"
	"saas/kernel/data"
	"saas/kernel/response"
	"strconv"
)

func DoTypeByCreate(ctx *gin.Context) {

	var former basic.DoTypeByCreateFormer
	if err := ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	tx := data.Database.Begin()

	typ := model.DorType{
		Name:     former.Name,
		Order:    former.Order,
		IsEnable: former.IsEnable,
	}

	if t := tx.Create(&typ); t.RowsAffected <= 0 {
		tx.Rollback()
		response.ToResponseByFail(ctx, "添加失败")
		return
	}

	if len(former.Beds) > 0 {
		beds := make([]model.DorTypeBed, len(former.Beds))
		for index, item := range former.Beds {
			beds[index] = model.DorTypeBed{TypeId: typ.Id, Name: item.Name, IsPublic: item.IsPublic}
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

func DoTypeByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.ToResponseByFailRequestMessage(ctx, "ID获取失败")
		return
	}

	var former basic.DoTypeByUpdateFormer
	if err = ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var typing model.DorType
	data.Database.First(&typing, id)
	if typing.Id <= 0 {
		response.ToResponseByNotFound(ctx, "未找到该房型")
		return
	}

	if typing.IsEnable != former.IsEnable {
		var peoples int64 = 0
		data.Database.Model(&model.DorPeople{}).Where("`type_id`=? and `status`=?", typing.Id, model.DorPeopleStatusLive).Count(&peoples)
		if peoples > 0 {
			response.ToResponseByFail(ctx, "该房型已有人入住，无法上下架")
			return
		}
	}

	typing.Name = former.Name
	typing.Order = former.Order
	typing.IsEnable = former.IsEnable

	if t := data.Database.Save(&typing); t.RowsAffected <= 0 {
		response.ToResponseByFail(ctx, "修改失败")
		return
	}

	response.ToResponseBySuccess(ctx)
}

func DoTypeByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.ToResponseByFailRequestMessage(ctx, "ID获取失败")
		return
	}

	var typing model.DorType
	data.Database.First(&typing, id)
	if typing.Id <= 0 {
		response.ToResponseByNotFound(ctx, "未找到该房型")
		return
	}

	var peoples int64 = 0
	data.Database.Model(model.DorPeople{}).Where("`type_id`=? and `status`=?", typing.Id, model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		response.ToResponseByFail(ctx, "该房型已有人入住，无法删除")
		return
	}

	tx := data.Database.Begin()

	if t := tx.Delete(&typing); t.RowsAffected <= 0 {
		tx.Rollback()
		response.ToResponseByFail(ctx, "删除失败")
		return
	}

	if t := tx.Where("type_id=?", typing.Id).Delete(&model.DorTypeBed{}); t.Error != nil {
		tx.Rollback()
		response.ToResponseByFail(ctx, "删除失败")
		return
	}

	tx.Commit()

	response.ToResponseBySuccess(ctx)
}

func DoTypeByEnable(ctx *gin.Context) {

	var former basic.DoTypeByEnableFormer
	if err := ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var typ model.DorType
	data.Database.First(&typ, former.Id)
	if typ.Id <= 0 {
		response.ToResponseByNotFound(ctx, "未找到该房型")
		return
	}

	var peoples int64 = 0
	data.Database.Model(&model.DorPeople{}).Where("`type_id`=? and `status`=?", typ.Id, model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		response.ToResponseByFail(ctx, "该房型已有人入住，无法上下架")
		return
	}

	typ.IsEnable = former.IsEnable

	if t := data.Database.Save(&typ); t.RowsAffected <= 0 {
		response.ToResponseByFail(ctx, "启禁失败")
		return
	}

	response.ToResponseBySuccess(ctx)
}

func ToTypeByList(ctx *gin.Context) {

	responses := make([]any, 0)

	var types []model.DorType
	data.Database.Preload("Beds").Order("`order` asc, `id` desc").Find(&types)

	for _, item := range types {
		items := basicResponse.ToTypeByListResponse{
			Id:        item.Id,
			Name:      item.Name,
			Order:     item.Order,
			IsEnable:  item.IsEnable,
			CreatedAt: item.CreatedAt.ToDateTimeString(),
		}
		for _, value := range item.Beds {
			items.Beds = append(items.Beds, basicResponse.ToTypeByListOfBedResponse{
				Name:     value.Name,
				IsPublic: value.IsPublic,
			})
		}
		responses = append(responses, items)
	}

	response.ToResponseBySuccessList(ctx, responses)
}

func ToTypeByOnline(ctx *gin.Context) {

	var query basic.ToTypeByOnlineFormer
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	responses := make([]any, 0)

	var types []model.DorType

	tx := data.Database.Where("`is_enable`=?", constant.IsEnableYes)

	if query.WithBed || query.MustBed {
		tx = tx.Preload("Beds")
	}
	if query.MustBed {
		tx = tx.Where("exists (?)", data.Database.
			Select("1").
			Table(model.TableDorTypeBed).
			Where(fmt.Sprintf("%s.id=%s.type_id", model.TableDorType, model.TableDorTypeBed)).
			Where(fmt.Sprintf("%s.deleted_at is null", model.TableDorTypeBed)),
		)
	}

	tx.Order("`order` asc, `id` desc").Find(&types)

	for _, item := range types {
		items := basicResponse.ToTypeByOnlineResponse{
			Id:   item.Id,
			Name: item.Name,
		}
		if query.WithBed || query.MustBed {
			for _, value := range item.Beds {
				items.Beds = append(items.Beds, basicResponse.ToTypeByOnlineOfBedResponse{
					Id:   value.Id,
					Name: value.Name,
				})
			}
		}
		responses = append(responses, items)
	}

	response.ToResponseBySuccessList(ctx, responses)
}
