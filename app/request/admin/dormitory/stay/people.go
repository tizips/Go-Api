package stay

import (
	"saas/app/request/basic"
)

type ToPeopleByPaginate struct {
	Building uint   `form:"building" json:"building" binding:"omitempty,gte=0"`
	Floor    uint   `form:"floor" json:"floor" binding:"omitempty,gte=0"`
	IsTemp   int8   `form:"is_temp" json:"is_temp" binding:"omitempty,eq=1|eq=2"`
	Status   string `form:"status" json:"status" binding:"required,eq=live|eq=leave"`
	Type     string `form:"type" json:"type" binding:"omitempty,eq=name|eq=mobile|eq=room"`
	Keyword  string `form:"keyword" json:"keyword" binding:"omitempty,max=20"`
	basic.Paginate
}

type DoPeopleByCreate struct {
	Bed      uint   `form:"bed" json:"bed" binding:"required,gt=0"`
	Category uint   `form:"category" json:"category" binding:"required,gt=0"`
	Name     string `form:"name" json:"name" binding:"required,max=20"`
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Start    string `form:"start" json:"start" binding:"required,datetime=2006-01-02"`
	End      string `form:"end" json:"end" binding:"omitempty,datetime,gtfield=Start"`
	Remark   string `form:"remark" json:"remark" binding:"omitempty,max=255"`
}
