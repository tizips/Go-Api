package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableSysRole = "sys_role"

type SysRole struct {
	Id        int `gorm:"primary_key"`
	Name      string
	Summary   string
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt

	BindPermissions []SysRoleBindPermission `gorm:"foreignKey:RoleId"`
}
