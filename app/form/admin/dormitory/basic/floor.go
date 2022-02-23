package basic

import "saas/app/form/basic"

type DoFloorByCreateForm struct {
	Name     string `form:"name" json:"name" binding:"required,max=20"`
	Building uint   `form:"building" json:"building" building:"required,gt=0"`
	Order    uint   `form:"order" json:"order" binding:"required,gte=1,lte=99"`
	basic.Enable
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
	Building uint `form:"building" json:"building" binding:"required,gt=0"`
}

type DoFloorByEnableForm struct {
	Id uint `form:"id" json:"id" binding:"required,gt=0"`
	basic.Enable
}
