package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableOaDepartmentOpen = "oa_department_open"

type OaDepartmentOpen struct {
	Id           uint `gorm:"primary_key"`
	DepartmentId uint
	Channel      string
	Openid       string
	CreatedAt    carbon.DateTime
	UpdatedAt    carbon.DateTime
	DeletedAt    gorm.DeletedAt
}
