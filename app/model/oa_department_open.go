package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableOaDepartmentOpen = "oa_department_open"

type OaDepartmentOpen struct {
	Id           int `gorm:"primary_key"`
	DepartmentId int
	Channel      string
	Openid       string
	CreatedAt    carbon.DateTime
	UpdatedAt    carbon.DateTime
	DeletedAt    gorm.DeletedAt
}
