package model

import "gorm.io/gorm"

const TableSysRoleBindPermission = "sys_role_bind_permission"

type SysRoleBindPermission struct {
	Id           int `gorm:"primary_key"`
	RoleId       int
	PermissionId int
	DeletedAt    gorm.DeletedAt

	Permission SysPermission `gorm:"foreignKey:Id;references:PermissionId"`
}
