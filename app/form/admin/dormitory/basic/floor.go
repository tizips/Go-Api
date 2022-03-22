package basic

import "saas/app/form/basic"

type DoFloorByCreateForm struct {
	Name     string `form:"name" json:"name" binding:"required,max=20"`
	Building uint   `form:"building" json:"building" building:"required,gt=0"`
	Order    uint   `form:"order" json:"order" binding:"required,gte=1,lte=99"`
	basic.Enable
	IsPublic uint8 `form:"is_public" json:"is_public" binding:"eq=1|eq=2"`
}

type DoFloorByUpdateForm struct {
	Name     string `form:"name" json:"name" binding:"required,max=20"`
	Building uint   `form:"building" json:"building" building:"required,gt=0"`
	Order    uint   `form:"order" json:"order" binding:"required,gte=1,lte=99"`
	basic.Enable
}

type ToFloorByListForm struct {
	Building uint `form:"building" json:"building" binding:"required,gt=0"`
}

type ToFloorByOnlineForm struct {
	Building   uint  `form:"building" json:"building" binding:"required,gt=0"`
	IsPublic   uint8 `form:"is_public" json:"is_public" binding:"omitempty,eq=1|eq=2"`
	WithPublic bool  `form:"with_public" json:"with_public" binding:"omitempty"`
}

type DoFloorByEnableForm struct {
	Id uint `form:"id" json:"id" binding:"required,gt=0"`
	basic.Enable
}
