package model

import "gorm.io/gorm"

const TableSysAdminBindRole = "sys_admin_bind_role"

type SysAdminBindRole struct {
	Id        uint `gorm:"primary_key"`
	AdminId   uint
	RoleId    uint
	DeletedAt gorm.DeletedAt

	Role SysRole `gorm:"References:RoleId;foreignKey:Id"`
}
