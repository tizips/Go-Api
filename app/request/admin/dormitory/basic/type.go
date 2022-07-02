package basic

import "saas/app/request/basic"

type DoTypeByCreate struct {
	Name  string `form:"name" json:"name" binding:"required,max=20" label:"名称"`
	Order int    `form:"order" json:"order" binding:"required,gte=1,lte=99" label:"序号"`
	Beds  []struct {
		Name     string `form:"name" json:"name" binding:"required,max=20" label:"名称"`
		IsPublic int8   `form:"is_public" json:"is_public" binding:"required,eq=1|eq=2" label:"装修"`
	} `form:"beds" json:"beds" binding:"omitempty,dive" label:"床位"`
	basic.Enable
}

type DoTypeByUpdate struct {
	Name  string `form:"name" json:"name" binding:"required,max=20" label:"名称"`
	Order int    `form:"order" json:"order" binding:"required,gte=1,lte=99" label:"序号"`
	basic.Enable
}

type DoTypeByEnable struct {
	Id int `form:"id" json:"id" binding:"required,gt=0"`
	basic.Enable
}

type ToTypeByOnline struct {
	MustBed bool `form:"must_bed" json:"must_bed" binding:"omitempty"`
	WithBed bool `form:"with_bed" json:"with_bed" binding:"omitempty"`
}
