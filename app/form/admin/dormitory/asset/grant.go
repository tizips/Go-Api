package asset

import "saas/app/form/basic"

type DoGrantByCreateFormer struct {
	Object    string                            `form:"object" json:"object" binding:"required,eq=package|eq=device"`
	Package   uint                              `form:"package" json:"package" binding:"required_if=Object package,omitempty,gt=0"`
	Device    uint                              `form:"device" json:"device" binding:"required_if=Object device,omitempty,gt=0"`
	Number    uint                              `form:"number" json:"number" binding:"required_if=Object device,omitempty,gt=0"`
	Position  string                            `form:"position" json:"position" binding:"required,eq=positions|eq=type"`
	Positions []DoGrantByCreateOfPositionFormer `form:"positions" json:"positions" binding:"required_if=Position positions,omitempty,min=1,dive"`
	Type      uint                              `form:"type" json:"type" binding:"required_if=Position type,omitempty,gt=0"`
	Remark    string                            `form:"remark" json:"remark" binding:"required,max=255"`
}

type DoGrantByCreateOfPositionFormer struct {
	Object string `form:"object" json:"object" binding:"required,eq=building|eq=floor|eq=room|eq=bed"`
	Id     uint   `form:"id" json:"id" binding:"required,gt=0"`
}

type DoGrantByRevokeFormer struct {
	Id uint `form:"id" json:"id" binding:"required,gt=0"`
}

type ToGrantByPaginateFormer struct {
	basic.Paginate
}
