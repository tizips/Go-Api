package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableSysModule = "sys_module"

type SysModule struct {
	Id        int `gorm:"primary_key"`
	Slug      string
	Name      string
	IsEnable  int8
	Order     int
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt

	SysPermissions []SysPermission `gorm:"foreignKey:ModuleId"`
}
