package auth

import "saas/app/request/basic"

type ToRoleByPaginate struct {
	basic.Paginate
}

type DoRoleByCreate struct {
	Name        string  `form:"name" json:"name" binding:"required,min=2,max=20" label:"名称"`
	Permissions [][]int `form:"permissions" json:"permissions" binding:"required" label:"权限"`
	Summary     string  `form:"summary" json:"summary" binding:"omitempty,max=255" label:"简介"`
}

type DoRoleByUpdate struct {
	Name        string  `form:"name" json:"name" binding:"required,min=2,max=20" label:"名称"`
	Permissions [][]int `form:"permissions" json:"permissions" binding:"required" label:"权限"`
	Summary     string  `form:"summary" json:"summary" binding:"omitempty,max=255" label:"简介"`
}
