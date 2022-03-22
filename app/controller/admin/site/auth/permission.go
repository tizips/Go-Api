package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"saas/app/constant"
	"saas/app/form/admin/site/auth"
	"saas/app/model"
	authResponse "saas/app/response/admin/site/auth"
	authService "saas/app/service/site/auth"
	authorize "saas/kernel/auth"
	"saas/kernel/data"
	"saas/kernel/response"
	"strconv"
)

func DoPermissionByCreate(ctx *gin.Context) {

	var former auth.DoPermissionByCreateForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var module model.SysModule
	data.Database.Where("is_enable=?", constant.IsEnableYes).First(&module, former.Module)
	if module.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "模块不存在",
		})
		return
	}

	var parent1, parent2 uint

	var parent model.SysPermission

	if former.Parent > 0 {
		data.Database.First(&parent, former.Parent)
		if parent.Id <= 0 {
			ctx.JSON(http.StatusOK, response.Response{
				Code:    40400,
				Message: "父级权限不存在",
			})
			return
		} else if parent.ParentI2 > 0 {
			ctx.JSON(http.StatusOK, response.Response{
				Code:    60000,
				Message: "该父级已是最低等级，无法继续添加",
			})
			return
		} else if parent.ParentI1 > 0 {

			if former.Method == "" || former.Path == "" {
				ctx.JSON(http.StatusOK, response.Response{
					Code:    60000,
					Message: "接口不能为空",
				})
				return
			}

			parent2 = parent.Id
			parent1 = parent.ParentI1
		} else {
			parent1 = parent.Id
		}
	}

	var permission model.SysPermission

	if former.Method != "" && former.Path != "" {
		var count int64
		data.Database.Model(model.SysPermission{}).Where("method = ?", former.Method).Where("path = ?", former.Path).Count(&count)
		if count > 0 {
			ctx.JSON(http.StatusOK, response.Response{
				Code:    60000,
				Message: "权限已存在",
			})
			return
		}
	}

	permission = model.SysPermission{
		ModuleId: module.Id,
		ParentI1: parent1,
		ParentI2: parent2,
		Name:     former.Name,
		Slug:     former.Slug,
		Method:   former.Method,
		Path:     former.Path,
	}

	data.Database.Create(&permission)
	if permission.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "添加失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})

}

func DoPermissionByUpdate(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))
	if id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "权限ID不存在",
		})
		return
	}

	var former auth.DoPermissionByUpdateForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var permission model.SysPermission
	data.Database.First(&permission, id)
	if permission.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "权限不存在",
		})
		return
	}

	method := permission.Method
	path := permission.Path

	var module model.SysModule
	data.Database.Where("is_enable=?", constant.IsEnableYes).First(&module, former.Module)
	if module.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "模块不存在",
		})
		return
	}

	var parent1, parent2 uint

	var parent model.SysPermission

	if former.Parent > 0 {
		data.Database.First(&parent, former.Parent)
		if parent.Id <= 0 {
			ctx.JSON(http.StatusOK, response.Response{
				Code:    40400,
				Message: "父级权限不存在",
			})
			return
		} else if parent.ParentI2 > 0 {
			ctx.JSON(http.StatusOK, response.Response{
				Code:    60000,
				Message: "该父级已是最低等级，无法继续添加",
			})
			return
		} else if parent.ParentI1 > 0 {

			if former.Method == "" || former.Path == "" {
				ctx.JSON(http.StatusOK, response.Response{
					Code:    60000,
					Message: "接口不能为空",
				})
				return
			}

			parent2 = parent.Id
			parent1 = parent.ParentI1
		} else {
			parent1 = parent.Id
		}
	}

	if former.Method != "" && former.Path != "" {
		var count int64
		data.Database.Model(model.SysPermission{}).Where("id <> ?", id).Where("method = ?", former.Method).Where("path = ?", former.Path).Count(&count)
		if count > 0 {
			ctx.JSON(http.StatusOK, response.Response{
				Code:    60000,
				Message: "权限已存在",
			})
			return
		}
	}

	permission.ParentI1 = parent1
	permission.ParentI2 = parent2
	permission.Name = former.Name
	permission.Slug = former.Slug
	permission.Method = former.Method
	permission.Path = former.Path

	tx := data.Database.Begin()

	if t := data.Database.Save(&permission); t.RowsAffected <= 0 {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "修改失败",
		})
		return
	}

	if method != former.Method || path != former.Path { //	变更权限
		if method != "" || path != "" {
			if _, err := authorize.Casbin.DeletePermission(method, path); err != nil {
				tx.Rollback()
				ctx.JSON(http.StatusOK, response.Response{
					Code:    60000,
					Message: "修改失败",
				})
				return
			}
		}

		if former.Method != "" || former.Path != "" {
			var bindings []model.SysRoleBindPermission
			tx.Where("permission_id = ?", permission.Id).Find(&bindings)
			if len(bindings) > 0 {
				for _, item := range bindings {
					if _, err := authorize.Casbin.AddPermissionForUser(authorize.NameByRole(item.RoleId), permission.Method, permission.Path); err != nil {
						tx.Rollback()
						ctx.JSON(http.StatusOK, response.Response{
							Code:    60000,
							Message: "修改失败",
						})
						return
					}
				}
			}
		}

	}

	tx.Commit()

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})

}

func DoPermissionByDelete(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))
	if id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "权限ID不存在",
		})
		return
	}

	var permission model.SysPermission
	data.Database.First(&permission, id)
	if permission.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "权限不存在",
		})
		return
	}

	tx := data.Database.Begin()

	if t := data.Database.Delete(&permission); t.RowsAffected <= 0 {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "模块删除失败",
		})
		return
	}

	if _, err := authorize.Casbin.DeletePermission(permission.Method, permission.Path); err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "模块删除失败",
		})
		return
	}

	tx.Commit()

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})

}

func ToPermissionByTree(ctx *gin.Context) {

	var former auth.ToPermissionByTreeForm
	if err := ctx.ShouldBindQuery(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Responses{
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

	results := authService.TreePermission(former.Module, false, false)
	for _, item := range results {
		responses.Data = append(responses.Data, item)
	}

	ctx.JSON(http.StatusOK, responses)
}

func ToPermissionByParents(ctx *gin.Context) {

	var former auth.ToPermissionByTreeForm
	if err := ctx.BindQuery(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Responses{
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

	results := authService.TreePermission(former.Module, true, true)
	for _, item := range results {
		responses.Data = append(responses.Data, item)
	}

	ctx.JSON(http.StatusOK, responses)
}

func ToPermissionBySelf(ctx *gin.Context) {

	responses := response.Responses{
		Code:    20000,
		Message: "Success",
		Data:    []any{},
	}

	var results []authResponse.TreePermissionResponse

	var modules []authResponse.TreePermissionResponse
	var children, children1, children2 []model.SysPermission

	if authorize.Root(authorize.Id(ctx)) {

		var permissions []model.SysPermission
		data.Database.Preload("Module").Find(&permissions)

		for _, item := range permissions {
			mark := true
			for _, value := range modules {
				if item.Module.Id == value.Id {
					mark = false
				}
			}
			if mark {
				modules = append(modules, authResponse.TreePermissionResponse{
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
			child := authResponse.TreePermissionResponse{
				Id:   item.Id,
				Name: item.Name,
			}

			for _, value := range children {
				if child.Id == value.Module.Id {

					//	处理第二层
					child1 := authResponse.TreePermissionResponse{
						Id:   value.Id,
						Name: value.Name,
					}

					for _, val := range children1 {
						if child1.Id == val.ParentI1 {

							//	处理第三层
							child2 := authResponse.TreePermissionResponse{
								Id:   val.Id,
								Name: val.Name,
							}

							for _, v := range children2 {
								if child2.Id == v.ParentI2 {

									//	处理第四层
									child3 := authResponse.TreePermissionResponse{
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
			responses.Data = append(responses.Data, item)
		}
	}

	ctx.JSON(http.StatusOK, responses)
}
