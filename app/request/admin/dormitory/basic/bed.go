package basic

import "saas/app/request/basic"

type DoBedByCreate struct {
	Name  string `form:"name" json:"name" binding:"required,max=20" label:"名称"`
	Room  int    `form:"room" json:"room" building:"required,gt=0" label:"房间"`
	Order int    `form:"order" json:"order" binding:"required,gte=1,lte=99" label:"序号"`
	basic.Enable
	IsPublic int8 `form:"is_public" json:"is_public" binding:"required,eq=1|eq=2" label:"公共"`
}

type DoBedByUpdate struct {
	Name  string `form:"name" json:"name" binding:"required,max=20" label:"名称"`
	Order int    `form:"order" json:"order" binding:"required,gte=1,lte=99" label:"序号"`
	basic.Enable
}

type ToBedByPaginate struct {
	Building int    `form:"building" json:"building" binding:"omitempty,gt=0" label:"楼栋"`
	Floor    int    `form:"floor" json:"floor" binding:"omitempty,gt=0" label:"楼层"`
	Room     int    `form:"room" json:"room" binding:"omitempty,gt=0" label:"房间"`
	Bed      string `form:"bed" json:"bed" binding:"omitempty,max=20" label:"床位"`
	basic.Paginate
	IsPublic int8 `form:"is_public" json:"is_public" binding:"omitempty,eq=1|eq=2" label:"公共"`
}

type ToBedByOnline struct {
	Room       int  `form:"room" json:"room" binding:"required,gt=0" label:"房间"`
	IsPublic   int8 `form:"is_public" json:"is_public" binding:"omitempty,eq=1|eq=2"`
	WithPublic bool `form:"with_public" json:"with_public" binding:"omitempty"`
}

type DoBedByEnable struct {
	Id int `form:"id" json:"id" binding:"required,gt=0" label:"ID"`
	basic.Enable
}
