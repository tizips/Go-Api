package asset

import "saas/app/request/basic"

type DoDeviceByCreate struct {
	Category      uint   `form:"category" json:"category" binding:"required,gt=0"`
	No            string `form:"no" json:"no" binding:"required,max=64"`
	Name          string `form:"name" json:"name" binding:"required,max=64"`
	Specification string `form:"specification" json:"specification" binding:"omitempty,max=64"`
	Stock         uint   `form:"stock" json:"stock" binding:"omitempty,gte=0"`
	Price         uint   `form:"price" json:"price" binding:"omitempty,gte=0"`
	Unit          string `form:"unit" json:"unit" binding:"required,max=64"`
	Indemnity     uint   `form:"indemnity" json:"indemnity" binding:"omitempty,gt=0"`
	Remark        string `form:"remark" json:"remark" binding:"omitempty,max=255"`
}

type DoDeviceByUpdate struct {
	Category      uint   `form:"category" json:"category" binding:"required,gt=0"`
	No            string `form:"no" json:"no" binding:"required,max=64"`
	Name          string `form:"name" json:"name" binding:"required,max=64"`
	Specification string `form:"specification" json:"specification" binding:"omitempty,max=64"`
	Stock         uint   `form:"stock" json:"stock" binding:"omitempty,gte=0"`
	Price         uint   `form:"price" json:"price" binding:"omitempty,gte=0"`
	Unit          string `form:"unit" json:"unit" binding:"required,max=64"`
	Indemnity     uint   `form:"indemnity" json:"indemnity" binding:"omitempty,gt=0"`
	Remark        string `form:"remark" json:"remark" binding:"omitempty,max=255"`
}

type ToDeviceByPaginate struct {
	basic.Paginate
	Type    string `form:"type" json:"type" binding:"omitempty,oneof=name no"`
	Keyword string `form:"keyword" json:"keyword" binding:"omitempty,max=64"`
}

type ToDeviceByOnline struct {
	Category uint `form:"category" json:"category" binding:"required,gt=0"`
}
