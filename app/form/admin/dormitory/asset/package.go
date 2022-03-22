package asset

import "saas/app/form/basic"

type DoPackageByCreateForm struct {
	Name    string `form:"name" json:"name" binding:"required,max=20"`
	Devices []struct {
		Device uint `form:"device" json:"device" binding:"required,gt=0"`
		Number uint `form:"number" json:"number" binding:"required,gt=0"`
	} `form:"devices" json:"devices" binding:"required,min=1,dive"`
}

type DoPackageByUpdateForm struct {
	Name    string `form:"name" json:"name" binding:"required,max=20"`
	Devices []struct {
		Device uint `form:"device" json:"device" binding:"required,gt=0"`
		Number uint `form:"number" json:"number" binding:"required,gt=0"`
	} `form:"devices" json:"devices" binding:"required,min=1,dive"`
}

type ToPackageByPaginateForm struct {
	basic.Paginate
	Keyword string `form:"keyword" json:"keyword" binding:"omitempty,max=20"`
}
