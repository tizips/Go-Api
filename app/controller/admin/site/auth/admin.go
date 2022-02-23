package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"saas/app/constant"
	"saas/app/form/admin/site/auth"
	"saas/app/model"
	authResponse "saas/app/response/admin/site/auth"
	authorize "saas/kernel/auth"
	"saas/kernel/data"
	"saas/kernel/response"
	"strconv"
)

func DoAdminByCreate(ctx *gin.Context) {

	var former auth.DoAdminByCreateForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var count int64

	tc := data.Database.Model(model.SysRole{})
	if !authorize.Root(authorize.Id(ctx)) {
		tc = tc.Where("id <> ?", authorize.ROOT)
	}
	tc.Where("id IN (?)", former.Roles).Count(&count)

	if len(former.Roles) != int(count) {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "部分角色不存在",
		})
		return
	}

	data.Database.Model(model.SysAdmin{}).Where("mobile = ?", former.Mobile).Count(&count)
	if count > 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "该手机号已被注册",
		})
		return
	}

	data.Database.Model(model.SysAdmin{}).Where("username = ?", former.Username).Count(&count)
	if count > 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "该用户名已被注册",
		})
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(former.Password), bcrypt.DefaultCost)

	tx := data.Database.Begin()

	admin := model.SysAdmin{
		Username: former.Username,
		Nickname: former.Nickname,
		Mobile:   former.Mobile,
		Password: string(password),
		IsEnable: former.IsEnable,
	}

	if tx.Create(&admin); admin.Id <= 0 {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "添加失败",
		})
		return
	}

	var binds []model.SysAdminBindRole

	for _, item := range former.Roles {
		binds = append(binds, model.SysAdminBindRole{
			AdminId: admin.Id,
			RoleId:  item,
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

	if len(binds) > 0 {
		var items = make([]string, len(binds))
		for idx, item := range binds {
			items[idx] = authorize.NameByRole(item.RoleId)
		}

		if _, err := authorize.Casbin.AddRolesForUser(authorize.NameByAdmin(admin.Id), items); err != nil {
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

func DoAdminByUpdate(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))
	if id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "账号ID不存在",
		})
		return
	}

	var former auth.DoAdminByUpdateForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var count int64
	tc := data.Database.Model(model.SysRole{})
	if !authorize.Root(authorize.Id(ctx)) {
		tc = tc.Where("id <> ?", authorize.ROOT)
	}
	tc.Where("id in (?)", former.Roles).Count(&count)

	if len(former.Roles) != int(count) {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "部分角色不存在",
		})
		return
	}

	data.Database.Model(model.SysAdmin{}).Where("id <> ?", id).Where("mobile = ?", former.Mobile).Count(&count)
	if count > 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "该手机号已被注册",
		})
		return
	}

	var admin model.SysAdmin

	data.Database.First(&admin, id)
	if admin.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "该账号不存在",
		})
		return
	}

	admin.Nickname = former.Nickname
	admin.Mobile = former.Mobile
	admin.IsEnable = former.IsEnable

	if former.Password != "" {

		password, _ := bcrypt.GenerateFromPassword([]byte(former.Password), bcrypt.DefaultCost)

		admin.Password = string(password)
	}

	var binds []model.SysAdminBindRole
	data.Database.Where("admin_id = ?", admin.Id).Find(&binds)

	var creates []model.SysAdminBindRole
	var deletes []uint
	var del []uint

	for _, item := range former.Roles {
		mark := true
		for _, value := range binds {
			if item == value.RoleId {
				mark = false
				break
			}
		}
		if mark {
			creates = append(creates, model.SysAdminBindRole{
				AdminId: admin.Id,
				RoleId:  item,
			})
		}
	}
	for _, item := range binds {
		mark := true
		for _, value := range former.Roles {
			if item.RoleId == value {
				mark = false
				break
			}
		}
		if mark {
			del = append(del, item.RoleId)
			deletes = append(deletes, item.Id)
		}
	}

	tx := data.Database.Begin()

	if t := tx.Save(&admin); t.RowsAffected <= 0 {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "修改失败",
		})
		return
	}

	if former.IsEnable == constant.IsEnableYes { //	用户禁用，删除缓存角色
		if _, err := authorize.Casbin.DeleteRolesForUser(authorize.NameByAdmin(admin.Id)); err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusOK, response.Response{
				Code:    60000,
				Message: "添加失败",
			})
			return
		}
	}

	if len(deletes) > 0 {
		var SysAdminBindRole model.SysAdminBindRole
		if t := tx.Where("admin_id = ?", admin.Id).Delete(&SysAdminBindRole, deletes); t.RowsAffected <= 0 {
			tx.Rollback()
			ctx.JSON(http.StatusOK, response.Response{
				Code:    60000,
				Message: "修改失败",
			})
			return
		}

		if len(del) > 0 && former.IsEnable == constant.IsEnableYes { //	用户启用，结算需要删除的角色
			for _, item := range del {
				if _, err := authorize.Casbin.DeleteRoleForUser(authorize.NameByAdmin(admin.Id), authorize.NameByRole(item)); err != nil {
					tx.Rollback()
					ctx.JSON(http.StatusOK, response.Response{
						Code:    60000,
						Message: "添加失败",
					})
					return
				}
			}
		}
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

		if len(creates) > 0 && former.IsEnable == constant.IsEnableYes { //	用户启用，处理需要新加的角色
			var items = make([]string, len(creates))
			for idx, item := range creates {
				items[idx] = authorize.NameByRole(item.RoleId)
			}

			if _, err := authorize.Casbin.AddRolesForUser(authorize.NameByAdmin(admin.Id), items); err != nil {
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

func ToAdminByPaginate(ctx *gin.Context) {

	var former auth.ToAdminByPaginateForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Paginate{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	tx := data.Database

	if !authorize.Root(authorize.Id(ctx)) {
		tx = tx.Where("not exists (?)", data.Database.
			Select("1").
			Model(model.SysAdminBindRole{}).
			Where(fmt.Sprintf("%s.id=%s.admin_id", model.TableSysAdmin, model.TableSysAdminBindRole)).
			Where("role_id = ?", authorize.ROOT),
		)
	}

	responses := response.Paginate{
		Code:    20000,
		Message: "Success",
	}

	responses.Data.Page = former.GetPage()
	responses.Data.Size = former.GetSize()
	responses.Data.Data = []interface{}{}

	tc := tx

	tc.Model(model.SysAdmin{}).Count(&responses.Data.Total)

	if responses.Data.Total > 0 {
		var admins []model.SysAdmin

		tx.
			Preload("BindRoles").
			Preload("BindRoles.Role").
			Order("id desc").
			Offset(former.GetOffset()).
			Limit(former.GetLimit()).
			Find(&admins)

		for _, item := range admins {
			result := authResponse.ToAdminByPaginateResponse{
				Id:        item.Id,
				Username:  item.Username,
				Nickname:  item.Nickname,
				Mobile:    item.Mobile,
				IsEnable:  item.IsEnable,
				CreatedAt: item.CreatedAt.ToDateTimeString(),
			}
			for _, value := range item.BindRoles {
				result.Roles = append(result.Roles, struct {
					Id   uint   `json:"id"`
					Name string `json:"name"`
				}{Id: value.Role.Id, Name: value.Role.Name})
			}
			responses.Data.Data = append(responses.Data.Data, result)
		}
	}

	ctx.JSON(http.StatusOK, responses)

}

func DoAdminByDelete(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))
	if id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "账号ID不存在",
		})
		return
	}

	if authorize.Id(ctx) == uint(id) {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "无法删除自身账号",
		})
		return
	}

	var admin model.SysAdmin
	data.Database.First(&admin, id)
	if admin.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "账号不存在",
		})
		return
	}

	tx := data.Database.Begin()

	if t := data.Database.Delete(&admin); t.RowsAffected <= 0 {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "账号删除失败",
		})
		return
	}

	bind := model.SysAdminBindRole{AdminId: admin.Id}

	if t := tx.Where("admin_id = ?", admin.Id).Delete(&bind); t.RowsAffected <= 0 {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "账号删除失败",
		})
		return
	}

	if _, err := authorize.Casbin.DeleteRolesForUser(authorize.NameByAdmin(admin.Id)); err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "添加失败",
		})
		return
	}

	tx.Commit()

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})

}

func DoAdminByEnable(ctx *gin.Context) {

	var former auth.DoAdminByEnableForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var admin model.SysAdmin
	data.Database.First(&admin, former.Id)
	if admin.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "账号不存在",
		})
		return
	}

	admin.IsEnable = former.IsEnable

	tx := data.Database.Begin()

	if t := data.Database.Save(&admin); t.RowsAffected <= 0 {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "启禁失败",
		})
		return
	}

	if former.IsEnable == constant.IsEnableNo {
		if _, err := authorize.Casbin.DeleteRolesForUser(authorize.NameByAdmin(admin.Id)); err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusOK, response.Response{
				Code:    60000,
				Message: "添加失败",
			})
			return
		}
	} else if former.IsEnable == constant.IsEnableYes {
		var roles []model.SysAdminBindRole
		tx.Where("admin_id = ?", admin.Id).Find(&roles)
		if len(roles) > 0 {
			var items []string
			for _, item := range roles {
				items = append(items, authorize.NameByRole(item.RoleId))
			}
			if _, err := authorize.Casbin.AddRolesForUser(authorize.NameByAdmin(admin.Id), items); err != nil {
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

	ctx.JSON(http.StatusOK, response.Responses{
		Code:    20000,
		Message: "Success",
	})

}
