package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableDorType = "dor_type"

type DorType struct {
	Id        int `gorm:"primary_key"`
	Name      string
	Order     int
	IsEnable  int8
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt

	Beds []DorTypeBed `gorm:"foreignKey:TypeId;references:Id"`
}
