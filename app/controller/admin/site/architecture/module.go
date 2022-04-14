package architecture

import (
	"github.com/gin-gonic/gin"
	"saas/app/constant"
	architectureForm "saas/app/form/admin/site/architecture"
	"saas/app/model"
	"saas/app/response/admin/site/architecture"
	"saas/kernel/data"
	"saas/kernel/response"
	"strconv"
)

func DoModuleByCreate(ctx *gin.Context) {

	var former architectureForm.DoModuleByCreateForm
	if err := ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var module model.SysModule
	data.Database.Where("slug = ?", former.Slug).First(&module)
	if module.Id > 0 {
		response.ToResponseByFail(ctx, "模块已存在")
		return
	}

	module = model.SysModule{
		Slug:     former.Slug,
		Name:     former.Name,
		IsEnable: former.IsEnable,
		Order:    former.Order,
	}

	data.Database.Create(&module)
	if module.Id <= 0 {
		response.ToResponseByFail(ctx, "模块创建失败")
		return
	}

	response.ToResponseBySuccess(ctx)
}

func DoModuleByUpdate(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))
	if id <= 0 {
		response.ToResponseByFailRequestMessage(ctx, "ID不存在")
		return
	}

	var former architectureForm.DoModuleByUpdateForm
	if err := ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var count int64
	data.Database.Model(model.SysModule{}).Where("`id`<>? and `slug`=?", id, former.Slug).Count(&count)
	if count > 0 {
		response.ToResponseByFail(ctx, "模块已存在")
		return
	}

	var module model.SysModule
	data.Database.First(&module, id)
	if module.Id <= 0 {
		response.ToResponseByNotFound(ctx, "模块不存在")
		return
	}

	module.Slug = former.Slug
	module.Name = former.Name
	module.IsEnable = former.IsEnable
	module.Order = former.Order

	tx := data.Database.Save(&module)
	if tx.RowsAffected <= 0 {
		response.ToResponseByFail(ctx, "模块修改失败")
		return
	}

	response.ToResponseBySuccess(ctx)
}

func DoModuleByDelete(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))
	if id <= 0 {
		response.ToResponseByFailRequestMessage(ctx, "ID不存在")
		return
	}

	var module model.SysModule
	data.Database.First(&module, id)
	if module.Id <= 0 {
		response.ToResponseByNotFound(ctx, "模块不存在")
		return
	}

	tx := data.Database.Delete(&module)
	if tx.RowsAffected <= 0 {
		response.ToResponseByFail(ctx, "模块删除失败")
		return
	}

	response.ToResponseBySuccess(ctx)
}

func DoModuleByEnable(ctx *gin.Context) {

	var former architectureForm.DoModuleByEnableForm
	if err := ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var module model.SysModule
	data.Database.First(&module, former.Id)
	if module.Id <= 0 {
		response.ToResponseByNotFound(ctx, "模块不存在")
		return
	}

	module.IsEnable = former.IsEnable

	tx := data.Database.Save(&module)
	if tx.RowsAffected <= 0 {
		response.ToResponseByFail(ctx, "启禁失败")
		return
	}

	response.ToResponseBySuccess(ctx)
}

func ToModuleByList(ctx *gin.Context) {

	responses := make([]any, 0)

	var modules []model.SysModule

	data.Database.Order("`order` asc").Find(&modules)

	for _, item := range modules {
		responses = append(responses, architecture.ToModuleByListResponse{
			Id:        item.Id,
			Slug:      item.Slug,
			Name:      item.Name,
			IsEnable:  item.IsEnable,
			Order:     item.Order,
			CreatedAt: item.CreatedAt.ToDateTimeString(),
		})
	}

	response.ToResponseBySuccessList(ctx, responses)
}

func ToModuleByOnline(ctx *gin.Context) {

	responses := make([]any, 0)

	var modules []model.SysModule

	data.Database.
		Where("`is_enable`=?", constant.IsEnableYes).
		Order("`order` asc").
		Find(&modules)

	for _, item := range modules {
		responses = append(responses, architecture.ToModuleByOnlineResponse{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	response.ToResponseBySuccessList(ctx, responses)
}
