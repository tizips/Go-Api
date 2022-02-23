package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableDorBuilding = "dor_building"

type DorBuilding struct {
	Id        uint `gorm:"primary_key"`
	Name      string
	Order     uint
	IsEnable  uint8
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt
}
