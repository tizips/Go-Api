package basic

import "saas/app/request/basic"

type DoBedByCreate struct {
	Name  string `form:"name" json:"name" binding:"required,max=20"`
	Room  uint   `form:"room" json:"room" building:"required,gt=0"`
	Order uint   `form:"order" json:"order" binding:"required,gte=1,lte=99"`
	basic.Enable
	IsPublic uint8 `form:"is_public" json:"is_public" binding:"required,eq=1|eq=2"`
}

type DoBedByUpdate struct {
	Name  string `form:"name" json:"name" binding:"required,max=20"`
	Order uint   `form:"order" json:"order" binding:"required,gte=1,lte=99"`
	basic.Enable
}

type ToBedByPaginate struct {
	Building uint   `form:"building" json:"building" binding:"omitempty,gt=0"`
	Floor    uint   `form:"floor" json:"floor" binding:"omitempty,gt=0"`
	Room     uint   `form:"room" json:"room" binding:"omitempty,gt=0"`
	Bed      string `form:"bed" json:"bed" binding:"omitempty,max=20"`
	basic.Paginate
	IsPublic uint8 `form:"is_public" json:"is_public" binding:"omitempty,eq=1|eq=2"`
}

type ToBedByOnline struct {
	Room       uint  `form:"room" json:"room" binding:"required,gt=0"`
	IsPublic   uint8 `form:"is_public" json:"is_public" binding:"omitempty,eq=1|eq=2"`
	WithPublic bool  `form:"with_public" json:"with_public" binding:"omitempty"`
}

type DoBedByEnable struct {
	Id uint `form:"id" json:"id" binding:"required,gt=0"`
	basic.Enable
}
