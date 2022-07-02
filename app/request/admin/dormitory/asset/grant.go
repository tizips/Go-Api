package asset

import "saas/app/request/basic"

type DoGrantByCreate struct {
	Object    string                      `form:"object" json:"object" binding:"required,eq=package|eq=device" label:"类型"`
	Package   int                         `form:"package" json:"package" binding:"required_if=Object package,omitempty,gt=0" label:"打包"`
	Device    int                         `form:"device" json:"device" binding:"required_if=Object device,omitempty,gt=0" label:"设备"`
	Number    int                         `form:"number" json:"number" binding:"required_if=Object device,omitempty,gt=0" label:"设备"`
	Position  string                      `form:"position" json:"position" binding:"required,eq=positions|eq=type" label:"位置"`
	Positions []DoGrantByCreateOfPosition `form:"positions" json:"positions" binding:"required_if=Position positions,omitempty,min=1,dive" label:"位置"`
	Type      int                         `form:"type" json:"type" binding:"required_if=Position type,omitempty,gt=0" label:"类别"`
	Remark    string                      `form:"remark" json:"remark" binding:"required,max=255" label:"备注"`
}

type DoGrantByCreateOfPosition struct {
	Object string `form:"object" json:"object" binding:"required,eq=building|eq=floor|eq=room|eq=bed" label:"类型"`
	Id     int    `form:"id" json:"id" binding:"required,gt=0" label:"ID"`
}

type DoGrantByRevoke struct {
	Id int `form:"id" json:"id" binding:"required,gt=0" label:"ID"`
}

type ToGrantByPaginate struct {
	basic.Paginate
}
