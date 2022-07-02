package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableOaDepartment = "oa_department"

type OaDepartment struct {
	Id        int `gorm:"primary_key"`
	ParentId  int
	Name      string
	IsEnable  int8
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt
}
