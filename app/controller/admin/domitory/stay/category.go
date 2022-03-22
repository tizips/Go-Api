package stay

import (
	"github.com/gin-gonic/gin"
	"net/http"
	liveForm "saas/app/form/admin/dormitory/stay"
	"saas/app/model"
	"saas/app/response/admin/dormitory/stay"
	"saas/kernel/data"
	"saas/kernel/response"
	"strconv"
)

func DoCategoryByCreate(ctx *gin.Context) {

	var former liveForm.DoCategoryByCreateForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	category := model.DorStayCategory{
		Name:     former.Name,
		Order:    former.Order,
		IsTemp:   former.IsTemp,
		IsEnable: former.IsEnable,
	}

	if data.Database.Create(&category); category.Id <= 0 {
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

func DoCategoryByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "类型ID获取失败",
		})
		return
	}

	var former liveForm.DoCategoryByUpdateForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var category model.DorStayCategory
	data.Database.First(&category, id)
	if category.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该类型",
		})
		return
	}

	if category.IsEnable != former.IsEnable {
		var peoples int64 = 0
		data.Database.Model(model.DorPeople{}).Where("category_id=?", category.Id).Where("status=?", model.DorPeopleStatusLive).Count(&peoples)
		if peoples > 0 {
			ctx.JSON(http.StatusOK, response.Response{
				Code:    40400,
				Message: "该类型正在使用，无法上下架",
			})
			return
		}
	}

	category.Name = former.Name
	category.Order = former.Order
	category.IsEnable = former.IsEnable

	if t := data.Database.Save(&category); t.RowsAffected <= 0 {
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

func DoCategoryByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "类型ID获取失败",
		})
		return
	}

	var category model.DorStayCategory
	data.Database.First(&category, id)
	if category.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该类型",
		})
		return
	}

	var peoples int64 = 0
	data.Database.Model(model.DorPeople{}).Where("category_id=?", category.Id).Where("status=?", model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "该类型正在使用，无法上下架",
		})
		return
	}

	if t := data.Database.Delete(&category); t.RowsAffected <= 0 {
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

func DoCategoryByEnable(ctx *gin.Context) {

	var former liveForm.DoCategoryByEnableForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var category model.DorStayCategory
	data.Database.First(&category, former.Id)
	if category.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该类型",
		})
		return
	}

	var peoples int64 = 0
	data.Database.Model(model.DorPeople{}).Where("category_id=?", category.Id).Where("status=?", model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "该类型正在使用，无法上下架",
		})
		return
	}

	category.IsEnable = former.IsEnable

	if t := data.Database.Save(&category); t.RowsAffected <= 0 {
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

func ToCategoryByList(ctx *gin.Context) {

	responses := response.Responses{
		Code:    20000,
		Message: "Success",
		Data:    []any{},
	}

	var categories []model.DorStayCategory
	data.Database.Order("`order` asc").Order("`id` desc").Find(&categories)

	for _, item := range categories {
		responses.Data = append(responses.Data, stay.ToCategoryByListResponse{
			Id:        item.Id,
			Name:      item.Name,
			Order:     item.Order,
			IsTemp:    item.IsTemp,
			IsEnable:  item.IsEnable,
			CreatedAt: item.CreatedAt.ToDateTimeString(),
		})
	}

	ctx.JSON(http.StatusOK, responses)
}

func ToCategoryByOnline(ctx *gin.Context) {

	responses := response.Responses{
		Code:    20000,
		Message: "Success",
		Data:    []any{},
	}

	var categories []model.DorStayCategory
	data.Database.Order("`order` asc").Order("`id` desc").Find(&categories)

	for _, item := range categories {
		responses.Data = append(responses.Data, stay.ToCategoryByOnlineResponse{
			Id:     item.Id,
			Name:   item.Name,
			IsTemp: item.IsTemp,
		})
	}

	ctx.JSON(http.StatusOK, responses)
}
