package basic

import "saas/app/request/basic"

type DoFloorByCreate struct {
	Name     string `form:"name" json:"name" binding:"required,max=20" label:"名称"`
	Building int    `form:"building" json:"building" building:"required,gt=0" label:"楼栋"`
	Order    int    `form:"order" json:"order" binding:"required,gte=1,lte=99" label:"序号"`
	basic.Enable
	IsPublic int8 `form:"is_public" json:"is_public" binding:"eq=1|eq=2" label:"公共"`
}

type DoFloorByUpdate struct {
	Name     string `form:"name" json:"name" binding:"required,max=20" label:"名称"`
	Building int    `form:"building" json:"building" building:"required,gt=0" label:"楼栋"`
	Order    int    `form:"order" json:"order" binding:"required,gte=1,lte=99" label:"序号"`
	basic.Enable
}

type ToFloorByList struct {
	Building int `form:"building" json:"building" binding:"required,gt=0" label:"楼栋"`
}

type ToFloorByOnline struct {
	Building   int  `form:"building" json:"building" binding:"required,gt=0" label:"楼栋"`
	IsPublic   int8 `form:"is_public" json:"is_public" binding:"omitempty,eq=1|eq=2"`
	WithPublic bool `form:"with_public" json:"with_public" binding:"omitempty"`
}

type DoFloorByEnable struct {
	Id int `form:"id" json:"id" binding:"required,gt=0"`
	basic.Enable
}
