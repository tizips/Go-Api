package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableOaDepartment = "oa_department"

type OaDepartment struct {
	Id        uint `gorm:"primary_key"`
	ParentId  uint
	Name      string
	IsEnable  uint8
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt
}
