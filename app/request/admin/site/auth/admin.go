package auth

import "saas/app/request/basic"

type ToAdminByPaginate struct {
	basic.Paginate
}

type DoAdminByCreate struct {
	Username string `form:"username" json:"username" binding:"required,min=4,max=20"`
	Nickname string `form:"nickname" json:"nickname" binding:"required,min=2,max=32"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=20"`
	Mobile   string `form:"mobile" json:"mobile" binding:"required,min=2,max=20"`
	Roles    []uint `form:"roles" json:"roles" binding:"required,unique,min=1"`
	basic.Enable
}

type DoAdminByUpdate struct {
	Nickname string `form:"nickname" json:"nickname" binding:"required,min=2,max=32"`
	Password string `form:"password" json:"password" binding:"omitempty,min=6,max=20"`
	Mobile   string `form:"mobile" json:"mobile" binding:"required,min=2,max=20"`
	Roles    []uint `form:"roles" json:"roles" binding:"required,unique,min=1"`
	basic.Enable
}

type DoAdminByEnable struct {
	Id uint `form:"id" json:"id" binding:"required,gt=0"`
	basic.Enable
}
