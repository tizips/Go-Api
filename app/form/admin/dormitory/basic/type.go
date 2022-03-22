package basic

import "saas/app/form/basic"

type DoTypeByCreateFormer struct {
	Name  string `form:"name" json:"name" binding:"required,max=20"`
	Order uint   `form:"order" json:"order" binding:"required,gte=1,lte=99"`
	Beds  []struct {
		Name     string `form:"name" json:"name" binding:"required,max=20"`
		IsPublic uint8  `form:"is_public" json:"is_public" binding:"required,eq=1|eq=2"`
	} `form:"beds" json:"beds" binding:"omitempty,dive"`
	basic.Enable
}

type DoTypeByUpdateFormer struct {
	Name  string `form:"name" json:"name" binding:"required,max=20"`
	Order uint   `form:"order" json:"order" binding:"required,gte=1,lte=99"`
	basic.Enable
}

type DoTypeByEnableFormer struct {
	Id uint `form:"id" json:"id" binding:"required,gt=0"`
	basic.Enable
}

type ToTypeByOnlineFormer struct {
	MustBed bool `form:"must_bed" json:"must_bed" binding:"omitempty"`
	WithBed bool `form:"with_bed" json:"with_bed" binding:"omitempty"`
}
