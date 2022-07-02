package basic

import "saas/app/request/basic"

type DoBuildingByCreate struct {
	Name  string `form:"name" json:"name" binding:"required,max=20" label:"名称"`
	Order int    `form:"order" json:"order" binding:"required,gte=1,lte=99" label:"序号"`
	basic.Enable
	IsPublic int8 `form:"is_public" json:"is_public" binding:"eq=1|eq=2" label:"公共"`
}

type DoBuildingByUpdate struct {
	Name  string `form:"name" json:"name" binding:"required,max=20" label:"名称"`
	Order int    `form:"order" json:"order" binding:"required,gte=1,lte=99" label:"序号"`
	basic.Enable
}

type DoBuildingByEnable struct {
	Id int `form:"id" json:"id" binding:"required,gt=0"`
	basic.Enable
}

type ToBuildingByOnline struct {
	IsPublic   int8 `form:"is_public" json:"is_public" binding:"omitempty,eq=1|eq=2"`
	WithPublic bool `form:"with_public" json:"with_public" binding:"omitempty"`
}
