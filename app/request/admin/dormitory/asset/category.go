package asset

import "saas/app/request/basic"

type DoCategoryByCreate struct {
	Name  string `form:"name" json:"name" binding:"required,max=20"`
	Order uint   `form:"order" json:"order" binding:"required,gte=1,lte=99"`
	basic.Enable
}

type DoCategoryByUpdate struct {
	Name  string `form:"name" json:"name" binding:"required,max=20"`
	Order uint   `form:"order" json:"order" binding:"required,gte=1,lte=99"`
	basic.Enable
}

type DoCategoryByEnable struct {
	Id uint `form:"id" json:"id" binding:"required,gt=0"`
	basic.Enable
}
