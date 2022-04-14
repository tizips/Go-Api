package stay

import (
	"github.com/gin-gonic/gin"
	"saas/app/constant"
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
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	category := model.DorStayCategory{
		Name:     former.Name,
		Order:    former.Order,
		IsTemp:   former.IsTemp,
		IsEnable: former.IsEnable,
	}

	if data.Database.Create(&category); category.Id <= 0 {
		response.ToResponseByFail(ctx, "添加失败")
		return
	}

	response.ToResponseBySuccess(ctx)
}

func DoCategoryByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.ToResponseByFailRequestMessage(ctx, "ID获取失败")
		return
	}

	var former liveForm.DoCategoryByUpdateForm
	if err = ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var category model.DorStayCategory
	data.Database.First(&category, id)
	if category.Id <= 0 {
		response.ToResponseByNotFound(ctx, "未找到该类型")
		return
	}

	if category.IsEnable != former.IsEnable {
		var peoples int64 = 0
		data.Database.Model(model.DorPeople{}).Where("`category_id`=? and `status`=?", category.Id, model.DorPeopleStatusLive).Count(&peoples)
		if peoples > 0 {
			response.ToResponseByFail(ctx, "该类型正在使用，无法上下架")
			return
		}
	}

	category.Name = former.Name
	category.Order = former.Order
	category.IsEnable = former.IsEnable

	if t := data.Database.Save(&category); t.RowsAffected <= 0 {
		response.ToResponseByFail(ctx, "修改失败")
		return
	}

	response.ToResponseBySuccess(ctx)
}

func DoCategoryByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.ToResponseByFailRequestMessage(ctx, "ID获取失败")
		return
	}

	var category model.DorStayCategory
	data.Database.First(&category, id)
	if category.Id <= 0 {
		response.ToResponseByNotFound(ctx, "未找到该类型")
		return
	}

	var peoples int64 = 0
	data.Database.Model(model.DorPeople{}).Where("`category_id`=? and `status`=?", category.Id, model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		response.ToResponseByFail(ctx, "该类型正在使用，无法上下架")
		return
	}

	if t := data.Database.Delete(&category); t.RowsAffected <= 0 {
		response.ToResponseByFail(ctx, "删除失败")
		return
	}

	response.ToResponseBySuccess(ctx)
}

func DoCategoryByEnable(ctx *gin.Context) {

	var former liveForm.DoCategoryByEnableForm
	if err := ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var category model.DorStayCategory
	data.Database.First(&category, former.Id)
	if category.Id <= 0 {
		response.ToResponseByNotFound(ctx, "未找到该类型")
		return
	}

	var peoples int64 = 0
	data.Database.Model(model.DorPeople{}).Where("`category_id`=? and `status`=?", category.Id, model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		response.ToResponseByFail(ctx, "该类型正在使用，无法上下架")
		return
	}

	category.IsEnable = former.IsEnable

	if t := data.Database.Save(&category); t.RowsAffected <= 0 {
		response.ToResponseByFail(ctx, "启禁失败")
		return
	}

	response.ToResponseBySuccess(ctx)
}

func ToCategoryByList(ctx *gin.Context) {

	responses := make([]any, 0)

	var categories []model.DorStayCategory
	data.Database.Order("`order` asc, `id` desc").Find(&categories)

	for _, item := range categories {
		responses = append(responses, stay.ToCategoryByListResponse{
			Id:        item.Id,
			Name:      item.Name,
			Order:     item.Order,
			IsTemp:    item.IsTemp,
			IsEnable:  item.IsEnable,
			CreatedAt: item.CreatedAt.ToDateTimeString(),
		})
	}

	response.ToResponseBySuccessList(ctx, responses)
}

func ToCategoryByOnline(ctx *gin.Context) {

	responses := make([]any, 0)

	var categories []model.DorStayCategory
	data.Database.Where("is_enable=?", constant.IsEnableYes).Order("`order` asc,`id` desc").Find(&categories)

	for _, item := range categories {
		responses = append(responses, stay.ToCategoryByOnlineResponse{
			Id:     item.Id,
			Name:   item.Name,
			IsTemp: item.IsTemp,
		})
	}

	response.ToResponseBySuccessList(ctx, responses)
}
