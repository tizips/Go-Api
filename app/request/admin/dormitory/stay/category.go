package stay

import "saas/app/request/basic"

type DoCategoryByCreate struct {
	Name   string `form:"name" json:"name" binding:"required,max=20" label:"名称"`
	Order  int    `form:"order" json:"order" binding:"required,gte=1,lte=99" label:"序号"`
	IsTemp int8   `form:"is_temp" json:"is_temp" binding:"oneof=0 1" label:"临时"`
	basic.Enable
}

type DoCategoryByUpdate struct {
	Name   string `form:"name" json:"name" binding:"required,max=20" label:"名称"`
	Order  int    `form:"order" json:"order" binding:"required,gte=1,lte=99" label:"序号"`
	IsTemp int8   `form:"is_temp" json:"is_temp" binding:"oneof=0 1" label:"临时"`
	basic.Enable
}

type DoCategoryByEnable struct {
	Id int `form:"id" json:"id" binding:"required,gt=0"`
	basic.Enable
}
