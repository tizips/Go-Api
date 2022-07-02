package basic

import "saas/app/request/basic"

type DoRoomByCreate struct {
	Name      string `form:"name" json:"name" binding:"required,max=20" label:"名称"`
	Floor     int    `form:"floor" json:"floor" building:"required,gt=0" label:"楼层"`
	Type      int    `form:"type" json:"type" building:"omitempty,required_if=IsPublic 2,gt=0" label:"类型"`
	Order     int    `form:"order" json:"order" binding:"required,gte=1,lte=99" label:"序号"`
	IsPublic  int8   `form:"is_public" json:"is_public" binding:"eq=1|eq=2" label:"公共"`
	IsFurnish int8   `form:"is_furnish" json:"is_furnish" binding:"omitempty,required_if=IsPublic 2,eq=1|eq=2" label:"装修"`
	basic.Enable
}

type DoRoomByUpdate struct {
	Name      string `form:"name" json:"name" binding:"required,max=20" label:"名称"`
	Type      int    `form:"type" json:"type" building:"omitempty,gt=0" label:"类型"`
	Order     int    `form:"order" json:"order" binding:"required,gte=1,lte=99" label:"序号"`
	IsFurnish int8   `form:"is_furnish" json:"is_furnish" binding:"omitempty,eq=1|eq=2" label:"装修"`
	basic.Enable
}

type ToRoomByPaginate struct {
	Building int    `form:"building" json:"building" binding:"omitempty,gt=0" label:"楼栋"`
	IsPublic int8   `form:"is_public" json:"is_public" binding:"omitempty,eq=1|eq=2" label:"公共"`
	Floor    int    `form:"floor" json:"floor" binding:"omitempty,gt=0" label:"楼层"`
	Room     string `form:"room" json:"room" binding:"omitempty,max=20" label:"房间"`
	basic.Paginate
}

type ToRoomByOnline struct {
	Floor      int  `form:"floor" json:"floor" binding:"required,gt=0"`
	IsPublic   int8 `form:"is_public" json:"is_public" binding:"omitempty,eq=1|eq=2"`
	WithPublic bool `form:"with_public" json:"with_public" binding:"omitempty"`
}

type DoRoomByEnable struct {
	Id int `form:"id" json:"id" binding:"required,gt=0"`
	basic.Enable
}

type DoRoomByFurnish struct {
	Id        int  `form:"id" json:"id" binding:"required,gt=0"`
	IsFurnish int8 `form:"is_furnish" json:"is_furnish" binding:"required,eq=1|eq=2"`
}
