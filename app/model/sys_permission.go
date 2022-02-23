package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableSysPermission = "sys_permission"

type SysPermission struct {
	Id        uint `gorm:"primary_key"`
	ModuleId  uint
	ParentI1  uint `gorm:"column:parent_i1"`
	ParentI2  uint `gorm:"column:parent_i2"`
	Name      string
	Slug      string
	Method    string
	Path      string
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt

	Module      SysModule       `gorm:"foreignKey:ModuleId"`
	Permissions []SysPermission `gorm:"many2many:sys_role_bind_permission;foreignKey:Id;joinForeignKey:RoleId;References:Id;joinReferences:PermissionId"`
}
