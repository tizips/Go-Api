package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableSysModule = "sys_module"

type SysModule struct {
	Id        uint `gorm:"primary_key"`
	Slug      string
	Name      string
	IsEnable  uint8
	Order     uint
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt

	SysPermissions []SysPermission `gorm:"foreignKey:ModuleId"`
}
