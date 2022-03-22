package basic

import "saas/app/form/basic"

type DoRoomByCreateFormer struct {
	Name      string `form:"name" json:"name" binding:"required,max=20"`
	Floor     uint   `form:"floor" json:"floor" building:"required,gt=0"`
	Type      uint   `form:"type" json:"type" building:"omitempty,required_if=IsPublic 2,gt=0"`
	Order     uint   `form:"order" json:"order" binding:"required,gte=1,lte=99"`
	IsPublic  uint8  `form:"is_public" json:"is_public" binding:"eq=1|eq=2"`
	IsFurnish uint8  `form:"is_furnish" json:"is_furnish" binding:"omitempty,required_if=IsPublic 2,eq=1|eq=2"`
	basic.Enable
}

type DoRoomByUpdateFormer struct {
	Name      string `form:"name" json:"name" binding:"required,max=20"`
	Type      uint   `form:"type" json:"type" building:"omitempty,gt=0"`
	Order     uint   `form:"order" json:"order" binding:"required,gte=1,lte=99"`
	IsFurnish uint8  `form:"is_furnish" json:"is_furnish" binding:"omitempty,eq=1|eq=2"`
	basic.Enable
}

type ToRoomByPaginateFormer struct {
	Building uint   `form:"building" json:"building" binding:"omitempty,gt=0"`
	IsPublic uint8  `form:"is_public" json:"is_public" binding:"omitempty,eq=1|eq=2"`
	Floor    uint   `form:"floor" json:"floor" binding:"omitempty,gt=0"`
	Room     string `form:"room" json:"room" binding:"omitempty,max=20"`
	basic.Paginate
}

type ToRoomByOnlineFormer struct {
	Floor      uint  `form:"floor" json:"floor" binding:"required,gt=0"`
	IsPublic   uint8 `form:"is_public" json:"is_public" binding:"omitempty,eq=1|eq=2"`
	WithPublic bool  `form:"with_public" json:"with_public" binding:"omitempty"`
}

type DoRoomByEnableFormer struct {
	Id uint `form:"id" json:"id" binding:"required,gt=0"`
	basic.Enable
}

type DoRoomByFurnishFormer struct {
	Id        uint  `form:"id" json:"id" binding:"required,gt=0"`
	IsFurnish uint8 `form:"is_furnish" json:"is_furnish" binding:"required,eq=1|eq=2"`
}
