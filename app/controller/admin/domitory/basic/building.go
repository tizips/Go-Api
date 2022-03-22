package basic

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	building := model.DorBuilding{
		Name:     former.Name,
		Order:    former.Order,
		IsEnable: former.IsEnable,
		IsPublic: former.IsPublic,
	}

	if data.Database.Create(&building); building.Id <= 0 {
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

func DoBuildingByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "楼栋ID获取失败",
		})
		return
	}

	var former basic.DoBuildingByUpdateFormer
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var building model.DorBuilding
	data.Database.First(&building, id)
	if building.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该楼栋",
		})
		return
	}

	if building.IsEnable != former.IsEnable {
		var peoples int64 = 0
		data.Database.Model(model.DorPeople{}).Where("building_id=?", building.Id).Where("status=?", model.DorPeopleStatusLive).Count(&peoples)
		if peoples > 0 {
			ctx.JSON(http.StatusOK, response.Response{
				Code:    40400,
				Message: "该楼栋已有人入住，无法上下架",
			})
			return
		}
	}

	building.Name = former.Name
	building.Order = former.Order
	building.IsEnable = former.IsEnable

	if t := data.Database.Save(&building); t.RowsAffected <= 0 {
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

func DoBuildingByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "楼栋ID获取失败",
		})
		return
	}

	var building model.DorBuilding
	data.Database.First(&building, id)
	if building.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该楼栋",
		})
		return
	}

	var peoples int64 = 0
	data.Database.Model(model.DorPeople{}).Where("building_id=?", building.Id).Where("status=?", model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "该楼栋已有人入住，无法删除",
		})
		return
	}

	tx := data.Database.Begin()

	if t := tx.Delete(&building); t.RowsAffected <= 0 {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "删除失败",
		})
		return
	}

	if t := tx.Where("building_id=?", building.Id).Delete(&model.DorRoom{}); t.Error != nil {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "删除失败",
		})
		return
	}

	if t := tx.Where("building_id=?", building.Id).Delete(&model.DorRoom{}); t.Error != nil {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "删除失败",
		})
		return
	}

	if t := tx.Where("building_id=?", building.Id).Delete(&model.DorBed{}); t.Error != nil {
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

func DoBuildingByEnable(ctx *gin.Context) {

	var former basic.DoBuildingByEnableFormer
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var building model.DorBuilding
	data.Database.First(&building, former.Id)
	if building.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该楼栋",
		})
		return
	}

	var peoples int64 = 0
	data.Database.Model(model.DorPeople{}).Where("building_id=?", building.Id).Where("status=?", model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "该楼栋已有人入住，无法上下架",
		})
		return
	}

	building.IsEnable = former.IsEnable

	if t := data.Database.Save(&building); t.RowsAffected <= 0 {
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

func ToBuildingByList(ctx *gin.Context) {

	responses := response.Responses{
		Code:    20000,
		Message: "Success",
		Data:    []any{},
	}

	var buildings []model.DorBuilding
	data.Database.Order("`order` asc").Order("`id` desc").Find(&buildings)

	for _, item := range buildings {
		responses.Data = append(responses.Data, basicResponse.ToBuildingByListResponse{
			Id:        item.Id,
			Name:      item.Name,
			Order:     item.Order,
			IsEnable:  item.IsEnable,
			IsPublic:  item.IsPublic,
			CreatedAt: item.CreatedAt.ToDateTimeString(),
		})
	}

	ctx.JSON(http.StatusOK, responses)
}

func ToBuildingByOnline(ctx *gin.Context) {

	var query basic.ToBuildingByOnlineFormer
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

	tx := data.Database

	if query.IsPublic > 0 {
		tx = tx.Where("is_public=?", query.IsPublic)
	}

	var buildings []model.DorBuilding
	tx.Order("`order` asc").Order("`id` desc").Find(&buildings)

	for _, item := range buildings {
		items := basicResponse.ToBuildingByOnlineResponse{
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
