package stay

import (
	"github.com/golang-module/carbon/v2"
	"saas/app/form/basic"
)

type ToPeopleByPaginateForm struct {
	Building uint   `form:"building" json:"building" binding:"omitempty,gte=0"`
	Floor    uint   `form:"floor" json:"floor" binding:"omitempty,gte=0"`
	IsTemp   int8   `form:"is_temp" json:"is_temp" binding:"omitempty,oneof=0 1"`
	Status   string `form:"status" json:"status" binding:"required,oneof=live leave"`
	Type     string `form:"type" json:"type" binding:"omitempty,oneof=name mobile room"`
	Keyword  string `form:"keyword" json:"keyword" binding:"omitempty,max=20"`
	basic.Paginate
}

type DoPeopleByCreateForm struct {
	Bed      uint        `form:"bed" json:"bed" binding:"required,gt=0"`
	Category uint        `form:"category" json:"category" binding:"required,gt=0"`
	Name     string      `form:"name" json:"name" binding:"required,max=20"`
	Mobile   string      `form:"mobile" json:"mobile" binding:"required,mobile"`
	Start    carbon.Date `form:"start" json:"start" binding:"required,datetime"`
	End      carbon.Date `form:"end" json:"end" binding:"omitempty,datetime,gtfield=Start"`
	Remark   string      `form:"remark" json:"remark" binding:"omitempty,max=255"`
}
