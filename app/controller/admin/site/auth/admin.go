package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var count int64

	tc := data.Database.Model(model.SysRole{})

	if !authorize.Root(authorize.Id(ctx)) {
		tc = tc.Where("`id`<>?", authorize.ROOT)
	}

	tc.Where("`id` IN (?)", former.Roles).Count(&count)

	if len(former.Roles) != int(count) {
		response.ToResponseByNotFound(ctx, "部分角色不存在")
		return
	}

	data.Database.Model(model.SysAdmin{}).Where("`mobile`=?", former.Mobile).Count(&count)
	if count > 0 {
		response.ToResponseByFail(ctx, "该手机号已被注册")
		return
	}

	data.Database.Model(model.SysAdmin{}).Where("username = ?", former.Username).Count(&count)
	if count > 0 {
		response.ToResponseByFail(ctx, "该用户名已被注册")
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
		response.ToResponseByFail(ctx, "添加失败")
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
		response.ToResponseByFail(ctx, "添加失败")
		return
	}

	if len(binds) > 0 {
		var items = make([]string, len(binds))
		for idx, item := range binds {
			items[idx] = authorize.NameByRole(item.RoleId)
		}

		if _, err := authorize.Casbin.AddRolesForUser(authorize.NameByAdmin(admin.Id), items); err != nil {
			tx.Rollback()
			response.ToResponseByFail(ctx, "添加失败")
			return
		}
	}

	tx.Commit()

	response.ToResponseBySuccess(ctx)
}

func DoAdminByUpdate(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))
	if id <= 0 {
		response.ToResponseByFailRequestMessage(ctx, "ID不存在")
		return
	}

	var former auth.DoAdminByUpdateForm
	if err := ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var count int64
	tc := data.Database.Model(model.SysRole{})
	if !authorize.Root(authorize.Id(ctx)) {
		tc = tc.Where("`id`<>?", authorize.ROOT)
	}
	tc.Where("`id` in (?)", former.Roles).Count(&count)

	if len(former.Roles) != int(count) {
		response.ToResponseByNotFound(ctx, "部分角色不存在")
		return
	}

	data.Database.Model(model.SysAdmin{}).Where("`id`<>? and `mobile`=?", id, former.Mobile).Count(&count)
	if count > 0 {
		response.ToResponseByFail(ctx, "该手机号已被注册")
		return
	}

	var admin model.SysAdmin

	data.Database.Preload("BindRoles").First(&admin, id)
	if admin.Id <= 0 {
		response.ToResponseByFail(ctx, "该账号不存在")
		return
	}

	admin.Nickname = former.Nickname
	admin.Mobile = former.Mobile
	admin.IsEnable = former.IsEnable

	if former.Password != "" {

		password, _ := bcrypt.GenerateFromPassword([]byte(former.Password), bcrypt.DefaultCost)

		admin.Password = string(password)
	}

	var creates []model.SysAdminBindRole
	var deletes []uint
	var del []uint

	for _, item := range former.Roles {
		mark := true
		for _, value := range admin.BindRoles {
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
	for _, item := range admin.BindRoles {
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
		response.ToResponseByFail(ctx, "修改失败")
		return
	}

	if former.IsEnable == constant.IsEnableYes { //	用户禁用，删除缓存角色
		if _, err := authorize.Casbin.DeleteRolesForUser(authorize.NameByAdmin(admin.Id)); err != nil {
			tx.Rollback()
			response.ToResponseByFail(ctx, "修改失败")
			return
		}
	}

	if len(deletes) > 0 {
		var SysAdminBindRole model.SysAdminBindRole
		if t := tx.Where("`admin_id`=?", admin.Id).Delete(&SysAdminBindRole, deletes); t.RowsAffected <= 0 {
			tx.Rollback()
			response.ToResponseByFail(ctx, "修改失败")
			return
		}

		if len(del) > 0 && former.IsEnable == constant.IsEnableYes { //	用户启用，结算需要删除的角色
			for _, item := range del {
				if _, err := authorize.Casbin.DeleteRoleForUser(authorize.NameByAdmin(admin.Id), authorize.NameByRole(item)); err != nil {
					tx.Rollback()
					response.ToResponseByFail(ctx, "修改失败")
					return
				}
			}
		}
	}

	if len(creates) > 0 {

		if t := tx.CreateInBatches(creates, 100); t.RowsAffected <= 0 {
			tx.Rollback()
			response.ToResponseByFail(ctx, "修改失败")
			return
		}

		if len(creates) > 0 && former.IsEnable == constant.IsEnableYes { //	用户启用，处理需要新加的角色
			var items = make([]string, len(creates))
			for idx, item := range creates {
				items[idx] = authorize.NameByRole(item.RoleId)
			}

			if _, err := authorize.Casbin.AddRolesForUser(authorize.NameByAdmin(admin.Id), items); err != nil {
				tx.Rollback()
				response.ToResponseByFail(ctx, "修改失败")
				return
			}
		}
	}

	tx.Commit()

	response.ToResponseBySuccess(ctx)
}

func ToAdminByPaginate(ctx *gin.Context) {

	var query auth.ToAdminByPaginateForm
	if err := ctx.ShouldBind(&query); err != nil {
		response.ToResponseByFailRequest(ctx, err)
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
		Total: 0,
		Page:  query.GetPage(),
		Size:  query.GetSize(),
		Data:  make([]any, 0),
	}

	tc := tx

	tc.Model(model.SysAdmin{}).Count(&responses.Total)

	if responses.Total > 0 {
		var admins []model.SysAdmin

		tx.
			//Preload("BindRoles").
			Preload("BindRoles.Role").
			Order("`id` desc").
			Offset(query.GetOffset()).
			Limit(query.GetLimit()).
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
			responses.Data = append(responses.Data, result)
		}
	}

	response.ToResponseBySuccessPaginate(ctx, responses)
}

func DoAdminByDelete(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))
	if id <= 0 {
		response.ToResponseByFailRequestMessage(ctx, "ID不存在")
		return
	}

	if authorize.Id(ctx) == uint(id) {
		response.ToResponseByFail(ctx, "无法删除自身账号")
		return
	}

	var admin model.SysAdmin
	data.Database.First(&admin, id)
	if admin.Id <= 0 {
		response.ToResponseByNotFound(ctx, "账号不存在")
		return
	}

	tx := data.Database.Begin()

	if t := data.Database.Delete(&admin); t.RowsAffected <= 0 {
		tx.Rollback()
		response.ToResponseByFail(ctx, "账号删除失败")
		return
	}

	bind := model.SysAdminBindRole{AdminId: admin.Id}

	if t := tx.Where("`admin_id`=?", admin.Id).Delete(&bind); t.RowsAffected <= 0 {
		tx.Rollback()
		response.ToResponseByFail(ctx, "账号删除失败")
		return
	}

	if _, err := authorize.Casbin.DeleteRolesForUser(authorize.NameByAdmin(admin.Id)); err != nil {
		tx.Rollback()
		response.ToResponseByFail(ctx, "账号删除失败")
		return
	}

	tx.Commit()

	response.ToResponseBySuccess(ctx)
}

func DoAdminByEnable(ctx *gin.Context) {

	var former auth.DoAdminByEnableForm
	if err := ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var admin model.SysAdmin
	data.Database.First(&admin, former.Id)
	if admin.Id <= 0 {
		response.ToResponseByNotFound(ctx, "账号不存在")
		return
	}

	admin.IsEnable = former.IsEnable

	tx := data.Database.Begin()

	if t := data.Database.Save(&admin); t.RowsAffected <= 0 {
		tx.Rollback()
		response.ToResponseByFail(ctx, "启禁失败")
		return
	}

	if former.IsEnable == constant.IsEnableNo {
		if _, err := authorize.Casbin.DeleteRolesForUser(authorize.NameByAdmin(admin.Id)); err != nil {
			tx.Rollback()
			response.ToResponseByFail(ctx, "启禁失败")
			return
		}
	} else if former.IsEnable == constant.IsEnableYes {
		tx.Where("`admin_id`=?", admin.Id).Find(&admin.BindRoles)
		if len(admin.BindRoles) > 0 {
			var items []string
			for _, item := range admin.BindRoles {
				items = append(items, authorize.NameByRole(item.RoleId))
			}
			if _, err := authorize.Casbin.AddRolesForUser(authorize.NameByAdmin(admin.Id), items); err != nil {
				tx.Rollback()
				response.ToResponseByFail(ctx, "启禁失败")
				return
			}
		}
	}

	tx.Commit()

	response.ToResponseBySuccess(ctx)
}
