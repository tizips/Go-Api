package asset

import "saas/app/request/basic"

type DoPackageByCreate struct {
	Name    string `form:"name" json:"name" binding:"required,max=20"`
	Devices []struct {
		Device uint `form:"device" json:"device" binding:"required,gt=0"`
		Number uint `form:"number" json:"number" binding:"required,gt=0"`
	} `form:"devices" json:"devices" binding:"required,min=1,unique=Device,dive"`
}

type DoPackageByUpdate struct {
	Name    string `form:"name" json:"name" binding:"required,max=20"`
	Devices []struct {
		Device uint `form:"device" json:"device" binding:"required,gt=0"`
		Number uint `form:"number" json:"number" binding:"required,gt=0"`
	} `form:"devices" json:"devices" binding:"required,min=1,unique=Device,dive"`
}

type ToPackageByPaginate struct {
	basic.Paginate
	Keyword string `form:"keyword" json:"keyword" binding:"omitempty,max=20"`
}
