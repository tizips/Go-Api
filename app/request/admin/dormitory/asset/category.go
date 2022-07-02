package asset

import "saas/app/request/basic"

type DoCategoryByCreate struct {
	Name  string `form:"name" json:"name" binding:"required,max=20" label:"名称"`
	Order int    `form:"order" json:"order" binding:"required,gte=1,lte=99" label:"序号"`
	basic.Enable
}

type DoCategoryByUpdate struct {
	Name  string `form:"name" json:"name" binding:"required,max=20" label:"名称"`
	Order int    `form:"order" json:"order" binding:"required,gte=1,lte=99" label:"序号"`
	basic.Enable
}

type DoCategoryByEnable struct {
	Id int `form:"id" json:"id" binding:"required,gt=0" label:"ID"`
	basic.Enable
}
