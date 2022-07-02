package asset

import "saas/app/request/basic"

type DoDeviceByCreate struct {
	Category      int    `form:"category" json:"category" binding:"required,gt=0" label:"类型"`
	No            string `form:"no" json:"no" binding:"required,max=64" label:"编号"`
	Name          string `form:"name" json:"name" binding:"required,max=64" label:"名称"`
	Specification string `form:"specification" json:"specification" binding:"omitempty,max=64" label:"规格"`
	Stock         int    `form:"stock" json:"stock" binding:"omitempty,gte=0" label:"库存"`
	Price         int    `form:"price" json:"price" binding:"omitempty,gte=0" label:"价格"`
	Unit          string `form:"unit" json:"unit" binding:"required,max=64" label:"单位"`
	Indemnity     int    `form:"indemnity" json:"indemnity" binding:"omitempty,gt=0" label:"赔偿金"`
	Remark        string `form:"remark" json:"remark" binding:"omitempty,max=255" label:"备注"`
}

type DoDeviceByUpdate struct {
	Category      int    `form:"category" json:"category" binding:"required,gt=0" label:"类型"`
	No            string `form:"no" json:"no" binding:"required,max=64" label:"编号"`
	Name          string `form:"name" json:"name" binding:"required,max=64" label:"名称"`
	Specification string `form:"specification" json:"specification" binding:"omitempty,max=64" label:"规格"`
	Stock         int    `form:"stock" json:"stock" binding:"omitempty,gte=0" label:"库存"`
	Price         int    `form:"price" json:"price" binding:"omitempty,gte=0" label:"价格"`
	Unit          string `form:"unit" json:"unit" binding:"required,max=64" label:"单位"`
	Indemnity     int    `form:"indemnity" json:"indemnity" binding:"omitempty,gt=0" label:"赔偿金"`
	Remark        string `form:"remark" json:"remark" binding:"omitempty,max=255" label:"备注"`
}

type ToDeviceByPaginate struct {
	basic.Paginate
	Type    string `form:"type" json:"type" binding:"omitempty,oneof=name no" label:"类别"`
	Keyword string `form:"keyword" json:"keyword" binding:"omitempty,max=64" label:"关键词"`
}

type ToDeviceByOnline struct {
	Category int `form:"category" json:"category" binding:"required,gt=0" label:"类型"`
}
