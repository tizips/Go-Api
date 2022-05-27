package basic

import "saas/app/request/basic"

type DoBuildingByCreate struct {
	Name  string `form:"name" json:"name" binding:"required,max=20"`
	Order uint   `form:"order" json:"order" binding:"required,gte=1,lte=99"`
	basic.Enable
	IsPublic uint8 `form:"is_public" json:"is_public" binding:"eq=1|eq=2"`
}

type DoBuildingByUpdate struct {
	Name  string `form:"name" json:"name" binding:"required,max=20"`
	Order uint   `form:"order" json:"order" binding:"required,gte=1,lte=99"`
	basic.Enable
}

type DoBuildingByEnable struct {
	Id uint `form:"id" json:"id" binding:"required,gt=0"`
	basic.Enable
}

type ToBuildingByOnline struct {
	IsPublic   uint8 `form:"is_public" json:"is_public" binding:"omitempty,eq=1|eq=2"`
	WithPublic bool  `form:"with_public" json:"with_public" binding:"omitempty"`
}
