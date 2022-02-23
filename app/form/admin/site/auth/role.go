package auth

import "saas/app/form/basic"

type ToRoleByPaginateForm struct {
	basic.Paginate
}

type DoRoleByCreateForm struct {
	Name        string   `form:"name" json:"name" binding:"required,min=2,max=20"`
	Permissions [][]uint `form:"permissions" json:"permissions" binding:"required"`
	Summary     string   `form:"summary" json:"summary" binding:"omitempty,max=255"`
}

type DoRoleByUpdateForm struct {
	Name        string   `form:"name" json:"name" binding:"required,min=2,max=20"`
	Permissions [][]uint `form:"permissions" json:"permissions" binding:"required"`
	Summary     string   `form:"summary" json:"summary" binding:"omitempty,max=255"`
}
