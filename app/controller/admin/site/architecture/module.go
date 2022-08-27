package architecture

import (
	"github.com/gin-gonic/gin"
	"saas/app/constant"
	"saas/app/model"
	architectureForm "saas/app/request/admin/site/architecture"
	"saas/app/response/admin/site/architecture"
	"saas/kernel/app"
	"saas/kernel/response"
	"strconv"
)

func DoModuleByCreate(ctx *gin.Context) {

	var request architectureForm.DoModuleByCreate
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var module model.SysModule
	app.MySQL.Where("slug = ?", request.Slug).Find(&module)
	if module.Id > 0 {
		response.Fail(ctx, "模块已存在")
		return
	}

	module = model.SysModule{
		Slug:     request.Slug,
		Name:     request.Name,
		IsEnable: request.IsEnable,
		Order:    request.Order,
	}

	app.MySQL.Create(&module)
	if module.Id <= 0 {
		response.Fail(ctx, "模块创建失败")
		return
	}

	response.Success(ctx)
}

func DoModuleByUpdate(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))
	if id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID不存在")
		return
	}

	var request architectureForm.DoModuleByUpdate
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var count int64
	app.MySQL.Model(model.SysModule{}).Where("`id`<>? and `slug`=?", id, request.Slug).Count(&count)
	if count > 0 {
		response.Fail(ctx, "模块已存在")
		return
	}

	var module model.SysModule
	app.MySQL.Find(&module, id)
	if module.Id <= 0 {
		response.NotFound(ctx, "模块不存在")
		return
	}

	module.Slug = request.Slug
	module.Name = request.Name
	module.IsEnable = request.IsEnable
	module.Order = request.Order

	tx := app.MySQL.Save(&module)
	if tx.RowsAffected <= 0 {
		response.Fail(ctx, "模块修改失败")
		return
	}

	response.Success(ctx)
}

func DoModuleByDelete(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))
	if id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID不存在")
		return
	}

	var module model.SysModule
	app.MySQL.Find(&module, id)
	if module.Id <= 0 {
		response.NotFound(ctx, "模块不存在")
		return
	}

	tx := app.MySQL.Delete(&module)
	if tx.RowsAffected <= 0 {
		response.Fail(ctx, "模块删除失败")
		return
	}

	response.Success(ctx)
}

func DoModuleByEnable(ctx *gin.Context) {

	var request architectureForm.DoModuleByEnable
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var module model.SysModule
	app.MySQL.Find(&module, request.Id)
	if module.Id <= 0 {
		response.NotFound(ctx, "模块不存在")
		return
	}

	module.IsEnable = request.IsEnable

	tx := app.MySQL.Save(&module)
	if tx.RowsAffected <= 0 {
		response.Fail(ctx, "启禁失败")
		return
	}

	response.Success(ctx)
}

func ToModuleByList(ctx *gin.Context) {

	responses := make([]any, 0)

	var modules []model.SysModule

	app.MySQL.Order("`order` asc").Find(&modules)

	for _, item := range modules {
		responses = append(responses, architecture.ToModuleByList{
			Id:        item.Id,
			Slug:      item.Slug,
			Name:      item.Name,
			IsEnable:  item.IsEnable,
			Order:     item.Order,
			CreatedAt: item.CreatedAt.ToDateTimeString(),
		})
	}

	response.SuccessByList(ctx, responses)
}

func ToModuleByOnline(ctx *gin.Context) {

	responses := make([]any, 0)

	var modules []model.SysModule

	app.MySQL.
		Where("`is_enable`=?", constant.IsEnableYes).
		Order("`order` asc").
		Find(&modules)

	for _, item := range modules {
		responses = append(responses, architecture.ToModuleByOnline{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	response.SuccessByList(ctx, responses)
}
