package manage

import (
	"github.com/gin-gonic/gin"
	"saas/app/constant"
	"saas/app/model"
	"saas/app/request/admin/site/manage"
	res "saas/app/response/admin/site/manage"
	authService "saas/app/service/site/manage"
	"saas/kernel/app"
	"saas/kernel/authorize"
	"saas/kernel/response"
	"strconv"
)

func DoPermissionByCreate(ctx *gin.Context) {

	var request manage.DoPermissionByCreate

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var module model.SysModule

	if app.MySQL.Where("`is_enable`=?", constant.IsEnableYes).Find(&module, request.Module); module.Id <= 0 {
		response.NotFound(ctx, "模块不存在")
		return
	}

	var parent1, parent2 int

	var parent model.SysPermission

	if request.Parent > 0 {

		if app.MySQL.Find(&parent, request.Parent); parent.Id <= 0 {
			response.Fail(ctx, "父级权限不存在")
			return
		} else if parent.ParentI2 > 0 {
			response.Fail(ctx, "该权限已是最低等级，无法继续添加")
			return
		} else if parent.ParentI1 > 0 {

			if request.Method == "" || request.Path == "" {
				response.Fail(ctx, "接口不能为空")
				return
			}

			parent2 = parent.Id
			parent1 = parent.ParentI1
		} else {
			parent1 = parent.Id
		}
	}

	var permission model.SysPermission

	if request.Method != "" && request.Path != "" {

		var count int64

		app.MySQL.Model(model.SysPermission{}).Where("`method`=? and `path`=?", request.Method, request.Path).Count(&count)

		if count > 0 {
			response.Fail(ctx, "权限已存在")
			return
		}
	}

	permission = model.SysPermission{
		ModuleId: module.Id,
		ParentI1: parent1,
		ParentI2: parent2,
		Name:     request.Name,
		Slug:     request.Slug,
		Method:   request.Method,
		Path:     request.Path,
	}

	if app.MySQL.Create(&permission); permission.Id <= 0 {
		response.Fail(ctx, "添加失败")
		return
	}

	response.Success(ctx)
}

func DoPermissionByUpdate(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))

	if id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID不存在")
		return
	}

	var request manage.DoPermissionByUpdate

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var permission model.SysPermission

	if app.MySQL.Find(&permission, id); permission.Id <= 0 {
		response.NotFound(ctx, "权限不存在")
		return
	}

	method := permission.Method
	path := permission.Path

	if permission.ModuleId != request.Module {

		var module model.SysModule

		if app.MySQL.Where("`is_enable`=?", constant.IsEnableYes).Find(&module, request.Module); module.Id <= 0 {
			response.NotFound(ctx, "模块不存在")
			return
		}
	}

	var parent1, parent2 int

	var parent model.SysPermission

	if request.Parent > 0 {

		if app.MySQL.Find(&parent, request.Parent); parent.Id <= 0 {
			response.NotFound(ctx, "父级权限不存在")
			return
		} else if parent.ParentI2 > 0 {
			response.Fail(ctx, "该权限已是最低等级，无法继续添加")
			return
		} else if parent.ParentI1 > 0 {

			if request.Method == "" || request.Path == "" {
				response.Fail(ctx, "接口不能为空")
				return
			}

			parent2 = parent.Id
			parent1 = parent.ParentI1
		} else {
			parent1 = parent.Id
		}
	}

	if request.Method != "" && request.Path != "" {

		var count int64

		app.MySQL.Model(model.SysPermission{}).Where("`id`<>? and `method`=? and `path`=?", id, request.Method, request.Path).Count(&count)

		if count > 0 {
			response.Fail(ctx, "权限已存在")
			return
		}
	}

	permission.ParentI1 = parent1
	permission.ParentI2 = parent2
	permission.Name = request.Name
	permission.Slug = request.Slug
	permission.Method = request.Method
	permission.Path = request.Path

	tx := app.MySQL.Begin()

	if t := app.MySQL.Save(&permission); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "修改失败")
		return
	}

	if method != request.Method || path != request.Path { //	变更权限
		if method != "" || path != "" {
			if _, err := app.Casbin.DeletePermission(method, path); err != nil {
				tx.Rollback()
				response.Fail(ctx, "修改失败")
				return
			}
		}

		if request.Method != "" || request.Path != "" {

			var bindings []model.SysRoleBindPermission

			if tx.Find(&bindings, "`permission_id`=?", permission.Id); len(bindings) > 0 {
				for _, item := range bindings {
					if _, err := app.Casbin.AddPermissionForUser(authorize.NameByRole(item.RoleId), permission.Method, permission.Path); err != nil {
						tx.Rollback()
						response.Fail(ctx, "修改失败")
						return
					}
				}
			}
		}

	}

	tx.Commit()

	response.Success(ctx)
}

func DoPermissionByDelete(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))

	if id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID不存在")
		return
	}

	var permission model.SysPermission

	if app.MySQL.Find(&permission, id); permission.Id <= 0 {
		response.Fail(ctx, "权限不存在")
		return
	}

	tx := app.MySQL.Begin()

	if t := tx.Delete(&permission); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "删除失败")
		return
	}

	if permission.Method != "" && permission.Path != "" {
		if _, err := app.Casbin.DeletePermission(permission.Method, permission.Path); err != nil {
			tx.Rollback()
			response.Fail(ctx, "删除失败")
			return
		}
	} else if permission.ParentI1 > 0 {

		var children []model.SysPermission

		if tx.Find(&children, "`parent_i2`=? and `method`<>? and `path`<>?", permission.Id, "", ""); len(children) > 0 {

			t := tx.Delete(&model.SysPermission{}, "`parent_i2`=?", permission.Id)

			if t.RowsAffected <= 0 {
				tx.Rollback()
				response.Fail(ctx, "删除失败")
				return
			}

			for _, item := range children {
				if _, err := app.Casbin.DeletePermission(item.Method, item.Path); err != nil {
					tx.Rollback()
					response.Fail(ctx, "删除失败")
					return
				}
			}
		}
	} else {

		var children []model.SysPermission

		if tx.Find(&children, "`parent_i1`=? and `method`<>? and `path`<>?", permission.Id, "", ""); len(children) > 0 {

			if t := tx.Delete(&model.SysPermission{}, "`parent_i1`=?", permission.Id); t.RowsAffected <= 0 {
				tx.Rollback()
				response.Fail(ctx, "删除失败")
				return
			}

			for _, item := range children {
				if _, err := app.Casbin.DeletePermission(item.Method, item.Path); err != nil {
					tx.Rollback()
					response.Fail(ctx, "删除失败")
					return
				}
			}
		}
	}

	tx.Commit()

	response.Success(ctx)
}

func ToPermissionByTree(ctx *gin.Context) {

	var request manage.ToPermissionByTree

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	responses := authService.TreePermission(request.Module, false, false)

	response.SuccessByData[[]res.TreePermission](ctx, responses)
}

func ToPermissionByParents(ctx *gin.Context) {

	var request manage.ToPermissionByTree

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	responses := authService.TreePermission(request.Module, true, true)

	response.SuccessByData[[]res.TreePermission](ctx, responses)
}

func ToPermissionBySelf(ctx *gin.Context) {

	responses := make([]any, 0)

	var results []res.TreePermission
	var modules []res.TreePermission

	var children, children1, children2 []model.SysPermission

	if authorize.Root(authorize.Id(ctx)) {

		var permissions []model.SysPermission

		app.MySQL.Preload("Module").Find(&permissions)

		for _, item := range permissions {

			mark := true

			for _, value := range modules {
				if item.Module.Id == value.Id {
					mark = false
				}
			}

			if mark {
				modules = append(modules, res.TreePermission{
					Id:   item.Module.Id,
					Name: item.Module.Name,
				})
			}
		}

		for _, item := range permissions {
			if item.ParentI2 > 0 {
				children2 = append(children2, item)
			} else if item.ParentI1 > 0 {
				children1 = append(children1, item)
			} else {
				children = append(children, item)
			}
		}
	}

	if len(modules) > 0 && (len(children2) > 0 || len(children1) > 0 || len(children) > 0) {

		for _, item := range modules {

			//	处理模块一层
			child := res.TreePermission{
				Id:   item.Id,
				Name: item.Name,
			}

			for _, value := range children {
				if child.Id == value.Module.Id {

					//	处理第二层
					child1 := res.TreePermission{
						Id:   value.Id,
						Name: value.Name,
					}

					for _, val := range children1 {
						if child1.Id == val.ParentI1 {

							//	处理第三层
							child2 := res.TreePermission{
								Id:   val.Id,
								Name: val.Name,
							}

							for _, v := range children2 {
								if child2.Id == v.ParentI2 {

									//	处理第四层
									child3 := res.TreePermission{
										Id:   v.Id,
										Name: v.Name,
									}

									if child3.Children != nil && len(child3.Children) > 0 || v.Method != "" && v.Path != "" {
										child2.Children = append(child2.Children, child3)
									}
								}
							}

							if child2.Children != nil && len(child2.Children) > 0 || val.Method != "" && val.Path != "" {
								child1.Children = append(child1.Children, child2)
							}
						}
					}

					if child1.Children != nil && len(child1.Children) > 0 || value.Method != "" && value.Path != "" {
						child.Children = append(child.Children, child1)
					}
				}
			}

			if child.Children != nil && len(child.Children) > 0 {
				results = append(results, child)
			}
		}

		for _, item := range results {
			responses = append(responses, item)
		}
	}

	response.SuccessByData(ctx, responses)
}
