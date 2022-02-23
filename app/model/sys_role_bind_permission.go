package model

import "gorm.io/gorm"

const TableSysRoleBindPermission = "sys_role_bind_permission"

type SysRoleBindPermission struct {
	Id           uint `gorm:"primary_key"`
	RoleId       uint
	PermissionId uint
	DeletedAt    gorm.DeletedAt

	Permission SysPermission `gorm:"foreignKey:Id;references:PermissionId"`
}
