package basic

import "saas/app/form/basic"

type DoRoomByCreateForm struct {
	Name      string `form:"name" json:"name" binding:"required,max=20"`
	Floor     uint   `form:"floor" json:"floor" building:"required,gt=0"`
	Type      uint   `form:"type" json:"type" building:"required,gt=0"`
	Order     uint   `form:"order" json:"order" binding:"required,gte=1,lte=99"`
	IsFurnish uint8  `form:"is_furnish" json:"is_furnish" binding:"oneof=0 1"`
	basic.Enable
}

type DoRoomByUpdateForm struct {
	Name      string `form:"name" json:"name" binding:"required,max=20"`
	Type      uint   `form:"type" json:"type" building:"required,gt=0"`
	Order     uint   `form:"order" json:"order" binding:"required,gte=1,lte=99"`
	IsFurnish uint8  `form:"is_furnish" json:"is_furnish" binding:"oneof=0 1"`
	basic.Enable
}

type ToRoomByPaginateForm struct {
	Building uint   `form:"building" json:"building" binding:"omitempty,gt=0"`
	Floor    uint   `form:"floor" json:"floor" binding:"omitempty,gt=0"`
	Room     string `form:"room" json:"room" binding:"omitempty,max=20"`
	basic.Paginate
}

type ToRoomByOnlineForm struct {
	Floor uint `form:"floor" json:"floor" binding:"required,gt=0"`
}

type DoRoomByEnableForm struct {
	Id uint `form:"id" json:"id" binding:"required,gt=0"`
	basic.Enable
}

type DoRoomByFurnishForm struct {
	Id        uint  `form:"id" json:"id" binding:"required,gt=0"`
	IsFurnish uint8 `form:"is_furnish" json:"is_furnish" binding:"oneof=0 1"`
}
