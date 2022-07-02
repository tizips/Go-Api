package auth

import "saas/app/request/basic"

type ToAdminByPaginate struct {
	basic.Paginate
}

type DoAdminByCreate struct {
	Username string `form:"username" json:"username" binding:"required,min=4,max=20" label:"用户名"`
	Nickname string `form:"nickname" json:"nickname" binding:"required,min=2,max=32" label:"昵称"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=20" label:"密码"`
	Mobile   string `form:"mobile" json:"mobile" binding:"required,min=2,max=20" label:"手机号"`
	Roles    []int  `form:"roles" json:"roles" binding:"required,unique,min=1" label:"角色"`
	basic.Enable
}

type DoAdminByUpdate struct {
	Nickname string `form:"nickname" json:"nickname" binding:"required,min=2,max=32" label:"昵称"`
	Password string `form:"password" json:"password" binding:"omitempty,min=6,max=20" label:"密码"`
	Mobile   string `form:"mobile" json:"mobile" binding:"required,min=2,max=20" label:"手机号"`
	Roles    []int  `form:"roles" json:"roles" binding:"required,unique,min=1" label:"角色"`
	basic.Enable
}

type DoAdminByEnable struct {
	Id int `form:"id" json:"id" binding:"required,gt=0"`
	basic.Enable
}
