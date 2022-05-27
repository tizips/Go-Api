package basic

import "saas/app/request/basic"

type DoTypeByCreate struct {
	Name  string `form:"name" json:"name" binding:"required,max=20"`
	Order uint   `form:"order" json:"order" binding:"required,gte=1,lte=99"`
	Beds  []struct {
		Name     string `form:"name" json:"name" binding:"required,max=20"`
		IsPublic uint8  `form:"is_public" json:"is_public" binding:"required,eq=1|eq=2"`
	} `form:"beds" json:"beds" binding:"omitempty,dive"`
	basic.Enable
}

type DoTypeByUpdate struct {
	Name  string `form:"name" json:"name" binding:"required,max=20"`
	Order uint   `form:"order" json:"order" binding:"required,gte=1,lte=99"`
	basic.Enable
}

type DoTypeByEnable struct {
	Id uint `form:"id" json:"id" binding:"required,gt=0"`
	basic.Enable
}

type ToTypeByOnline struct {
	MustBed bool `form:"must_bed" json:"must_bed" binding:"omitempty"`
	WithBed bool `form:"with_bed" json:"with_bed" binding:"omitempty"`
}
