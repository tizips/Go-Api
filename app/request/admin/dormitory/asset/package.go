package asset

import "saas/app/request/basic"

type DoPackageByCreate struct {
	Name    string `form:"name" json:"name" binding:"required,max=20" label:"名称"`
	Devices []struct {
		Device int `form:"device" json:"device" binding:"required,gt=0" label:"设备"`
		Number int `form:"number" json:"number" binding:"required,gt=0" label:"数量"`
	} `form:"devices" json:"devices" binding:"required,min=1,unique=Device,dive" label:"设备"`
}

type DoPackageByUpdate struct {
	Name    string `form:"name" json:"name" binding:"required,max=20" label:"名称"`
	Devices []struct {
		Device int `form:"device" json:"device" binding:"required,gt=0" label:"设备"`
		Number int `form:"number" json:"number" binding:"required,gt=0" label:"数量"`
	} `form:"devices" json:"devices" binding:"required,min=1,unique=Device,dive" label:"设备"`
}

type ToPackageByPaginate struct {
	basic.Paginate
	Keyword string `form:"keyword" json:"keyword" binding:"omitempty,max=20" label:"关键词"`
}
