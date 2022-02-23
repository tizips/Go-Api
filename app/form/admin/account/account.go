package account

type ToAccountByPermissionForm struct {
	Module uint `form:"module" binding:"required,numeric,gt=0"`
}
