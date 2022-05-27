package auth

import "saas/app/request/basic"

type ToRoleByPaginate struct {
	basic.Paginate
}

type DoRoleByCreate struct {
	Name        string   `form:"name" json:"name" binding:"required,min=2,max=20"`
	Permissions [][]uint `form:"permissions" json:"permissions" binding:"required"`
	Summary     string   `form:"summary" json:"summary" binding:"omitempty,max=255"`
}

type DoRoleByUpdate struct {
	Name        string   `form:"name" json:"name" binding:"required,min=2,max=20"`
	Permissions [][]uint `form:"permissions" json:"permissions" binding:"required"`
	Summary     string   `form:"summary" json:"summary" binding:"omitempty,max=255"`
}
