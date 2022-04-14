package basic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"saas/app/constant"
	accountForm "saas/app/form/admin/account"
	"saas/app/model"
	"saas/app/response/admin/account"
	"saas/kernel/auth"
	"saas/kernel/data"
	"saas/kernel/response"
)

func ToAccountByInformation(ctx *gin.Context) {

	admin := auth.Admin(ctx)

	response.ToResponseBySuccessData(ctx, account.ToAccountByInformationResponse{
		Username: admin.Username,
		Nickname: admin.Nickname,
		Avatar:   admin.Avatar,
		Mobile:   admin.Mobile,
	})
}

func ToAccountByModule(ctx *gin.Context) {

	responses := make([]any, 0)

	var modules []model.SysModule

	tx := data.Database.
		Where("`is_enable` = ?", constant.IsEnableYes)

	tc := data.Database.
		Select("1").
		Model(model.SysPermission{}).
		Where(fmt.Sprintf("`%s`.`id`=`%s`.`module_id`", model.TableSysModule, model.TableSysPermission))

	if !auth.Root(auth.Id(ctx)) {
		tc = tc.
			Joins(fmt.Sprintf("left join `%s` on `%s`.`id`=`%s`.`permission_id`", model.TableSysRoleBindPermission, model.TableSysPermission, model.TableSysRoleBindPermission)).
			Joins(fmt.Sprintf("left join `%s` on `%s`.`role_id`=`%s`.`role_id` and `%s`.`admin_id`=?", model.TableSysAdminBindRole, model.TableSysRoleBindPermission, model.TableSysAdminBindRole, model.TableSysAdminBindRole), auth.Id(ctx)).
			Where(fmt.Sprintf("`%s`.`id` is not null and `%s`.`deleted_at` is null and `%s`.`deleted_at` is null", model.TableSysAdminBindRole, model.TableSysRoleBindPermission, model.TableSysAdminBindRole))
	}

	tx.
		Where("exists (?)", tc).
		Order("`order` asc").
		Find(&modules)

	for _, item := range modules {
		responses = append(responses, account.ToAccountByModuleResponse{
			Id:   item.Id,
			Slug: item.Slug,
			Name: item.Name,
		})
	}

	response.ToResponseBySuccessList(ctx, responses)
}

func ToAccountByPermission(ctx *gin.Context) {

	var former accountForm.ToAccountByPermissionForm

	if err := ctx.BindQuery(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var responses = make([]any, 0)

	var permissions []model.SysPermission

	tx := data.Database.
		Where("`module_id` = ? and `method` <> ? and `path` <> ?", former.Module, "", "")

	if !auth.Root(auth.Id(ctx)) {
		tx = tx.
			Joins(fmt.Sprintf("left join `%s` on `%s`.`id`=`%s`.`permission_id`", model.TableSysRoleBindPermission, model.TableSysPermission, model.TableSysRoleBindPermission)).
			Joins(fmt.Sprintf("left join `%s` on `%s`.`role_id`=`%s`.`role_id` and `%s`.`admin_id`=?", model.TableSysAdminBindRole, model.TableSysRoleBindPermission, model.TableSysAdminBindRole, model.TableSysAdminBindRole), auth.Id(ctx)).
			Where(fmt.Sprintf("`%s`.`id` is not null and `%s`.`deleted_at` is null and `%s`.`deleted_at` is null", model.TableSysAdminBindRole, model.TableSysRoleBindPermission, model.TableSysAdminBindRole))
	}

	tx.
		Find(&permissions)

	for _, item := range permissions {
		responses = append(responses, item.Slug)
	}

	response.ToResponseBySuccessList(ctx, responses)
}
