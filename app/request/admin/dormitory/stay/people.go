package stay

import (
	"saas/app/request/basic"
)

type ToPeopleByPaginate struct {
	Building int    `form:"building" json:"building" binding:"omitempty,gte=0" label:"楼栋"`
	Floor    int    `form:"floor" json:"floor" binding:"omitempty,gte=0" label:"楼层"`
	IsTemp   int8   `form:"is_temp" json:"is_temp" binding:"omitempty,eq=1|eq=2" label:"临时"`
	Status   string `form:"status" json:"status" binding:"required,eq=live|eq=leave" label:"状态"`
	Type     string `form:"type" json:"type" binding:"omitempty,eq=name|eq=mobile|eq=room" label:"类型"`
	Keyword  string `form:"keyword" json:"keyword" binding:"omitempty,max=20" label:"关键词"`
	basic.Paginate
}

type DoPeopleByCreate struct {
	Bed      int    `form:"bed" json:"bed" binding:"required,gt=0" label:"床位"`
	Category int    `form:"category" json:"category" binding:"required,gt=0" label:"类型"`
	Name     string `form:"name" json:"name" binding:"required,max=20" label:"姓名"`
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile" label:"手机号"`
	Start    string `form:"start" json:"start" binding:"required,datetime=2006-01-02" label:"入住日期"`
	End      string `form:"end" json:"end" binding:"omitempty,datetime,gtfield=Start" label:"预离日期"`
	Remark   string `form:"remark" json:"remark" binding:"omitempty,max=255" label:"备注"`
}

type DoPeopleByLeave struct {
	Id     int    `json:"id" form:"id" binding:"required,gt=0"`
	Notify string `json:"notify" form:"notify" binding:"omitempty,max=255"`
	Remark string `json:"remark" form:"remark" binding:"required,max=255" label:"备注"`
}
