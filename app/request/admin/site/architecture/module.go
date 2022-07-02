package architecture

import "saas/app/request/basic"

type DoModuleByCreate struct {
	Slug  string `form:"slug" json:"slug" binding:"required,min=2,max=20,alpha" label:"标识"`
	Name  string `form:"name" json:"name" binding:"required,min=2,max=20" label:"名称"`
	Order int    `form:"order" json:"order" binding:"omitempty,gt=1,lt=99" label:"序号"`
	basic.Enable
}

type DoModuleByUpdate struct {
	Slug  string `form:"slug" json:"slug" binding:"required,min=2,max=20,alpha" label:"标识"`
	Name  string `form:"name" json:"name" binding:"required,min=2,max=20" label:"名称"`
	Order int    `form:"order" json:"order" binding:"omitempty,gt=1,lt=99" label:"序号"`
	basic.Enable
}

type DoModuleByEnable struct {
	Id int `form:"id" json:"id" binding:"required,gt=0"`
	basic.Enable
}
