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

func DoBedByCreate(ctx *gin.Context) {

	var former basic.DoBedByCreateFormer
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var room model.DorRoom
	data.Database.Where("is_enable=?", constant.IsEnableYes).First(&room, former.Room)
	if room.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "房间不存在",
		})
		return
	}
	if room.IsPublic == model.DorBedIsPublicYes {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "该房间为公共区域，无法添加",
		})
		return
	}

	bed := model.DorBed{
		BuildingId: room.BuildingId,
		FloorId:    room.FloorId,
		RoomId:     room.Id,
		TypeId:     room.TypeId,
		Name:       former.Name,
		Order:      former.Order,
		IsEnable:   former.IsEnable,
		IsPublic:   former.IsPublic,
	}

	if data.Database.Create(&bed); room.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "添加失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})

}

func DoBedByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "床位ID获取失败",
		})
		return
	}

	var former basic.DoBedByUpdateFormer
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var bed model.DorBed
	data.Database.First(&bed, id)
	if bed.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该床位",
		})
		return
	}

	if bed.IsEnable != former.IsEnable {
		var peoples int64 = 0
		data.Database.Model(model.DorPeople{}).Where("bed_id=?", bed.Id).Where("status=?", model.DorPeopleStatusLive).Count(&peoples)
		if peoples > 0 {
			ctx.JSON(http.StatusOK, response.Response{
				Code:    40400,
				Message: "该床位已有人入住，无法上下架",
			})
			return
		}
	}

	bed.Name = former.Name
	bed.Order = former.Order
	bed.IsEnable = former.IsEnable

	if t := data.Database.Save(&bed); t.RowsAffected <= 0 {
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

func DoBedByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "床位ID获取失败",
		})
		return
	}

	var bed model.DorBed
	data.Database.First(&bed, id)
	if bed.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该床位",
		})
		return
	}

	var peoples int64 = 0
	data.Database.Model(model.DorPeople{}).Where("bed_id=?", bed.Id).Where("status=?", model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "该床位已有人入住，无法删除",
		})
		return
	}

	if t := data.Database.Delete(&bed); t.RowsAffected <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "删除失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})

}

func DoBedByEnable(ctx *gin.Context) {

	var former basic.DoBedByEnableFormer
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var bed model.DorBed
	data.Database.First(&bed, former.Id)
	if bed.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该床位",
		})
		return
	}

	var peoples int64 = 0
	data.Database.Model(model.DorPeople{}).Where("bed_id=?", bed.Id).Where("status=?", model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "该床位已有人入住，无法上下架",
		})
		return
	}

	bed.IsEnable = former.IsEnable

	if t := data.Database.Save(&bed); t.RowsAffected <= 0 {
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

func ToBedByPaginate(ctx *gin.Context) {

	var query basic.ToBedByPaginateFormer
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	responses := response.Paginate{
		Code:    20000,
		Message: "Success",
	}

	responses.Data.Size = query.GetSize()
	responses.Data.Page = query.GetPage()
	responses.Data.Data = []any{}

	tx := data.Database

	if query.Room > 0 {
		tx = tx.Where("room_id", query.Floor)
	} else if query.Floor > 0 {
		tx = tx.Where("floor_id", query.Floor)
	} else if query.Building > 0 {
		tx = tx.Where("building_id", query.Building)
	}

	if query.Bed != "" {
		tx = tx.Where("name like ?", "%"+query.Bed+"%")
	}

	tc := tx

	tc.Model(&model.DorBed{}).Count(&responses.Data.Total)

	if responses.Data.Total > 0 {

		var beds []model.DorBed

		tx.
			Preload("Building").
			Preload("Floor").
			Preload("Room").
			Order("`order` asc").
			Order("`id` desc").
			Offset(query.GetOffset()).
			Limit(query.GetLimit()).
			Find(&beds)

		for _, item := range beds {
			responses.Data.Data = append(responses.Data.Data, basicResponse.ToBedByPaginateResponse{
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

	ctx.JSON(http.StatusOK, responses)
}

func ToBedByOnline(ctx *gin.Context) {

	var query basic.ToBedByOnlineFormer
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

	tx := data.Database.Where("room_id=?", query.Room)

	if query.IsPublic > 0 {
		tx = tx.Where("is_public=?", query.IsPublic)
	}

	var beds []model.DorBed
	tx.Order("`order` asc").Order("`id` desc").Find(&beds)

	for _, item := range beds {
		items := basicResponse.ToBedByOnlineResponse{
			Id:   item.Id,
			Name: item.Name,
		}
		if query.WithPublic {
			items.IsPublic = item.IsPublic
		}
		responses.Data = append(responses.Data, items)
	}

	ctx.JSON(http.StatusOK, responses)
}
