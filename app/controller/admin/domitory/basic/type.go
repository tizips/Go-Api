package basic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
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
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "添加失败",
		})
		return
	}

	if len(former.Beds) > 0 {
		beds := make([]model.DorTypeBed, len(former.Beds))
		for index, item := range former.Beds {
			beds[index] = model.DorTypeBed{TypeId: typ.Id, Name: item.Name, IsPublic: item.IsPublic}
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

func DoTypeByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "房型ID获取失败",
		})
		return
	}

	var former basic.DoTypeByUpdateFormer
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var typ model.DorType
	data.Database.First(&typ, id)
	if typ.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该房型",
		})
		return
	}

	if typ.IsEnable != former.IsEnable {
		var peoples int64 = 0
		data.Database.Model(model.DorPeople{}).Where("type_id=?", typ.Id).Where("status=?", model.DorPeopleStatusLive).Count(&peoples)
		if peoples > 0 {
			ctx.JSON(http.StatusOK, response.Response{
				Code:    40400,
				Message: "该房型已有人入住，无法上下架",
			})
			return
		}
	}

	typ.Name = former.Name
	typ.Order = former.Order
	typ.IsEnable = former.IsEnable

	if t := data.Database.Save(&typ); t.RowsAffected <= 0 {
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

func DoTypeByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "房型ID获取失败",
		})
		return
	}

	var typ model.DorType
	data.Database.First(&typ, id)
	if typ.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该房型",
		})
		return
	}

	var peoples int64 = 0
	data.Database.Model(model.DorPeople{}).Where("type_id=?", typ.Id).Where("status=?", model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "该房型已有人入住，无法删除",
		})
		return
	}

	tx := data.Database.Begin()

	if t := tx.Delete(&typ); t.RowsAffected <= 0 {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "删除失败",
		})
		return
	}

	if t := tx.Where("type_id=?", typ.Id).Delete(&model.DorTypeBed{}); t.Error != nil {
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

func DoTypeByEnable(ctx *gin.Context) {

	var former basic.DoTypeByEnableFormer
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var typ model.DorType
	data.Database.First(&typ, former.Id)
	if typ.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该房型",
		})
		return
	}

	var peoples int64 = 0
	data.Database.Model(model.DorPeople{}).Where("type_id=?", typ.Id).Where("status=?", model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "该房型已有人入住，无法上下架",
		})
		return
	}

	typ.IsEnable = former.IsEnable

	if t := data.Database.Save(&typ); t.RowsAffected <= 0 {
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

func ToTypeByList(ctx *gin.Context) {

	responses := response.Responses{
		Code:    20000,
		Message: "Success",
		Data:    []any{},
	}

	var types []model.DorType
	data.Database.Preload("Beds").Order("`order` asc").Order("`id` desc").Find(&types)

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
		responses.Data = append(responses.Data, items)
	}

	ctx.JSON(http.StatusOK, responses)
}

func ToTypeByOnline(ctx *gin.Context) {

	var query basic.ToTypeByOnlineFormer
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

	var types []model.DorType

	tx := data.Database
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
	tx.Order("`order` asc").Order("`id` desc").Find(&types)

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
		responses.Data = append(responses.Data, items)
	}

	ctx.JSON(http.StatusOK, responses)
}
