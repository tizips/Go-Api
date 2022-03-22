package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"saas/app/form/admin/site/auth"
	"saas/app/model"
	authResponse "saas/app/response/admin/site/auth"
	authorize "saas/kernel/auth"
	"saas/kernel/data"
	"saas/kernel/response"
	"strconv"
)

func DoRoleByCreate(ctx *gin.Context) {

	var former auth.DoRoleByCreateForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var permissionsIds []uint

	var modules, children1, children2, children3 []uint

	for _, item := range former.Permissions {
		if len(item) >= 4 {
			children3 = append(children3, item[3])
		} else if len(item) >= 3 {
			children2 = append(children2, item[2])
		} else if len(item) >= 2 {
			children1 = append(children1, item[1])
		} else if len(item) >= 1 {
			modules = append(modules, item[0])
		}
	}

	if len(modules) > 0 {
		var permissions []model.SysPermission
		data.Database.Where("module_id in (?)", modules).Where("method <> ?", "").Where("path <> ?", "").Find(&permissions)
		for _, item := range permissions {
			permissionsIds = append(permissionsIds, item.Id)
		}
	}
	if len(children3) > 0 {
		var permissions []model.SysPermission
		data.Database.Where("id in (?)", children3).Where("method <> ?", "").Where("path <> ?", "").Find(&permissions)
		for _, item := range permissions {
			permissionsIds = append(permissionsIds, item.Id)
		}
	}
	if len(children2) > 0 {
		var permissions []model.SysPermission
		data.Database.Where("parent_i2 in (?)", children2).Where("method <> ?", "").Where("path <> ?", "").Find(&permissions)
		for _, item := range permissions {
			permissionsIds = append(permissionsIds, item.Id)
		}
	}
	if len(children1) > 0 {
		var permissions []model.SysPermission
		data.Database.Where("parent_i1 in (?)", children1).Where("method <> ?", "").Where("path <> ?", "").Find(&permissions)
		for _, item := range permissions {
			permissionsIds = append(permissionsIds, item.Id)
		}
	}

	var temp = make(map[uint]uint, len(permissionsIds))
	for _, item := range permissionsIds {
		temp[item] = item
	}

	var bind []uint
	for _, item := range temp {
		bind = append(bind, item)
	}

	if len(bind) <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "可用权限不能为空",
		})
		return
	}

	tx := data.Database.Begin()

	role := model.SysRole{
		Name:    former.Name,
		Summary: former.Summary,
	}

	if tx.Create(&role); role.Id <= 0 {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "添加失败",
		})
		return
	}

	var binds []model.SysRoleBindPermission

	for _, item := range bind {
		binds = append(binds, model.SysRoleBindPermission{
			RoleId:       role.Id,
			PermissionId: item,
		})
	}

	if t := tx.CreateInBatches(binds, 100); t.RowsAffected <= 0 {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "添加失败",
		})
		return
	}

	var permissions []model.SysRoleBindPermission
	tx.
		Preload("Permission",
			tx.
				Where("method <> ?", "").
				Where("path <> ?", ""),
		).
		Where("role_id = ?", role.Id).
		Find(&permissions)

	if len(permissions) > 0 {
		var items = make([][]string, len(permissions))
		for idx, item := range permissions {
			items[idx] = []string{item.Permission.Method, item.Permission.Path}
		}

		if _, err := authorize.Casbin.AddPermissionsForUser(authorize.NameByRole(role.Id), items...); err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusOK, response.Response{
				Code:    60000,
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

func DoRoleByUpdate(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))
	if id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "角色ID不存在",
		})
		return
	}

	if id == authorize.ROOT {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "开发组权限，无法修改",
		})
		return
	}

	var former auth.DoRoleByUpdateForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var role model.SysRole

	data.Database.First(&role, id)
	if role.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "角色不存在",
		})
		return
	}

	var permissionsIds []uint

	var modules, children1, children2, children3 []uint

	for _, item := range former.Permissions {
		if len(item) >= 4 {
			children3 = append(children3, item[3])
		} else if len(item) >= 3 {
			children2 = append(children2, item[2])
		} else if len(item) >= 2 {
			children1 = append(children1, item[1])
		} else if len(item) >= 1 {
			modules = append(modules, item[0])
		}
	}

	if len(modules) > 0 {
		var permissions []model.SysPermission
		data.Database.Where("module_id in (?)", modules).Where("method <> ?", "").Where("path <> ?", "").Find(&permissions)
		for _, item := range permissions {
			permissionsIds = append(permissionsIds, item.Id)
		}
	}
	if len(children3) > 0 {
		var permissions []model.SysPermission
		data.Database.Where("id in (?)", children3).Where("method <> ?", "").Where("path <> ?", "").Find(&permissions)
		for _, item := range permissions {
			permissionsIds = append(permissionsIds, item.Id)
		}
	}
	if len(children2) > 0 {
		var permissions []model.SysPermission
		data.Database.Where("parent_i2 in (?)", children2).Where("method <> ?", "").Where("path <> ?", "").Find(&permissions)
		for _, item := range permissions {
			permissionsIds = append(permissionsIds, item.Id)
		}
	}
	if len(children1) > 0 {
		var permissions []model.SysPermission
		data.Database.Where("parent_i1 in (?)", children1).Where("method <> ?", "").Where("path <> ?", "").Find(&permissions)
		for _, item := range permissions {
			permissionsIds = append(permissionsIds, item.Id)
		}
	}

	var temp = make(map[uint]uint, len(permissionsIds))
	for _, item := range permissionsIds {
		temp[item] = item
	}

	var bind []uint
	for _, item := range temp {
		bind = append(bind, item)
	}

	if len(bind) <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "可用权限不能为空",
		})
		return
	}

	role.Name = former.Name
	role.Summary = former.Summary

	var binds []model.SysRoleBindPermission
	data.Database.Where("role_id = ?", role.Id).Find(&binds)

	var creates []model.SysRoleBindPermission
	var deletes []uint

	if len(binds) > 0 {
		for _, item := range bind {
			mark := true
			for _, value := range binds {
				if item == value.PermissionId {
					mark = false
					break
				}
			}
			if mark {
				creates = append(creates, model.SysRoleBindPermission{
					RoleId:       role.Id,
					PermissionId: item,
				})
			}
		}
		for _, item := range binds {
			mark := true
			for _, value := range bind {
				if item.PermissionId == value {
					mark = false
					break
				}
			}
			if mark {
				deletes = append(deletes, item.Id)
			}
		}
	}

	tx := data.Database.Begin()

	if t := tx.Save(&role); t.RowsAffected <= 0 {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "修改失败",
		})
		return
	}

	if len(creates) > 0 {

		if t := tx.CreateInBatches(creates, 100); t.RowsAffected <= 0 {
			tx.Rollback()
			ctx.JSON(http.StatusOK, response.Response{
				Code:    60000,
				Message: "修改失败",
			})
			return
		}

		var ids []uint
		for _, item := range creates {
			ids = append(ids, item.PermissionId)
		}
		var permissions []model.SysPermission
		tx.Where("method <> ?", "").Where("path <> ?", "").Find(&permissions, ids)
		if len(permissions) > 0 {
			var items [][]string
			for _, item := range permissions {
				items = append(items, []string{item.Method, item.Path})
			}
			if _, err := authorize.Casbin.AddPermissionsForUser(authorize.NameByRole(role.Id), items...); err != nil {
				tx.Rollback()
				ctx.JSON(http.StatusOK, response.Response{
					Code:    60000,
					Message: "修改失败",
				})
				return
			}
		}
	}

	if len(deletes) > 0 {
		var b model.SysRoleBindPermission
		if t := tx.Where("role_id = ?", role.Id).Delete(&b, deletes); t.RowsAffected <= 0 {
			tx.Rollback()
			ctx.JSON(http.StatusOK, response.Response{
				Code:    60000,
				Message: "修改失败",
			})
			return
		}
	}

	if len(deletes) > 0 {
		if _, err := authorize.Casbin.DeletePermissionsForUser(authorize.NameByRole(role.Id)); err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusOK, response.Response{
				Code:    60000,
				Message: "修改失败",
			})
			return
		}
	}

	if len(creates) > 0 || len(deletes) > 0 {
		var permissions []model.SysRoleBindPermission
		tx.
			Preload("Permission",
				tx.
					Where("method <> ?", "").
					Where("path <> ?", ""),
			).
			Where("role_id = ?", role.Id).
			Find(&permissions)

		if len(permissions) > 0 {
			var items = make([][]string, len(permissions))
			for idx, item := range permissions {
				items[idx] = []string{item.Permission.Method, item.Permission.Path}
			}

			if _, err := authorize.Casbin.AddPermissionsForUser(authorize.NameByRole(role.Id), items...); err != nil {
				tx.Rollback()
				ctx.JSON(http.StatusOK, response.Response{
					Code:    60000,
					Message: "添加失败",
				})
				return
			}
		}
	}

	tx.Commit()

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})

}

func DoRoleByDelete(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))
	if id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "角色ID不存在",
		})
		return
	}

	if id == authorize.ROOT {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "开发组权限，无法修改",
		})
		return
	}

	var role model.SysRole
	data.Database.Find(&role, id)
	if role.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "角色不存在",
		})
		return
	}

	tx := data.Database.Begin()

	if t := data.Database.Delete(&role); t.RowsAffected <= 0 {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "角色删除失败",
		})
		return
	}

	bind := model.SysRoleBindPermission{RoleId: role.Id}

	if t := tx.Where("role_id = ?", role.Id).Delete(&bind); t.RowsAffected <= 0 {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "角色删除失败",
		})
		return
	}

	if _, err := authorize.Casbin.DeleteRole(authorize.NameByRole(role.Id)); err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "角色删除失败",
		})
		return
	}

	tx.Commit()

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})

}

func ToRoleByPaginate(ctx *gin.Context) {

	var former auth.ToRoleByPaginateForm
	if err := ctx.ShouldBindQuery(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Responses{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	tx := data.Database.Where("id <> ?", authorize.ROOT)

	responses := response.Paginate{
		Code:    20000,
		Message: "Success",
	}

	responses.Data.Size = former.GetSize()
	responses.Data.Page = former.GetPage()
	responses.Data.Data = []any{}

	tc := tx

	tc.Model(model.SysRole{}).Count(&responses.Data.Total)

	if responses.Data.Total > 0 {

		tx = tx.Order("`id` desc")

		var roles []model.SysRole

		tx.Preload("BindPermissions").Preload("BindPermissions.Permission").Offset(former.GetOffset()).Limit(former.GetLimit()).Find(&roles)

		for _, item := range roles {

			r := authResponse.ToRoleByPaginateResponse{
				Id:        item.Id,
				Name:      item.Name,
				Summary:   item.Summary,
				CreatedAt: item.CreatedAt.ToDateTimeString(),
			}

			for _, value := range item.BindPermissions {
				r.Permissions = append(r.Permissions, []uint{value.Permission.ModuleId, value.Permission.ParentI1, value.Permission.ParentI2, value.PermissionId})
			}

			responses.Data.Data = append(responses.Data.Data, r)
		}
	}

	ctx.JSON(http.StatusOK, responses)
}

func ToRoleByEnable(ctx *gin.Context) {

	var roles []model.SysRole

	tx := data.Database

	if !authorize.Root(authorize.Id(ctx)) {
		tx.Where("role_id <> ?", authorize.ROOT)
	}

	tx.Find(&roles)

	responses := response.Responses{
		Code:    20000,
		Message: "Success",
	}

	for _, item := range roles {
		responses.Data = append(responses.Data, authResponse.ToRoleByEnableResponse{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	ctx.JSON(http.StatusOK, responses)

}
