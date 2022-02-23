package basic

import "saas/app/form/basic"

type DoBedByCreateForm struct {
	Name  string `form:"name" json:"name" binding:"required,max=20"`
	Room  uint   `form:"room" json:"room" building:"required,gt=0"`
	Order uint   `form:"order" json:"order" binding:"required,gte=1,lte=99"`
	basic.Enable
}

type DoBedByUpdateForm struct {
	Name  string `form:"name" json:"name" binding:"required,max=20"`
	Order uint   `form:"order" json:"order" binding:"required,gte=1,lte=99"`
	basic.Enable
}

type ToBedByPaginateForm struct {
	Building uint   `form:"building" json:"building" binding:"omitempty,gt=0"`
	Floor    uint   `form:"floor" json:"floor" binding:"omitempty,gt=0"`
	Room     uint   `form:"room" json:"room" binding:"omitempty,gt=0"`
	Bed      string `form:"bed" json:"bed" binding:"omitempty,max=20"`
	basic.Paginate
}

type ToBedByOnlineForm struct {
	Room uint `form:"room" json:"room" binding:"required,gt=0"`
}

type DoBedByEnableForm struct {
	Id uint `form:"id" json:"id" binding:"required,gt=0"`
	basic.Enable
}
