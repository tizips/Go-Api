package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableDorFloor = "dor_floor"

type DorFloor struct {
	Id         int `gorm:"primary_key"`
	BuildingId int
	Name       string
	Order      int
	IsEnable   int8
	IsPublic   int8
	CreatedAt  carbon.DateTime
	UpdatedAt  carbon.DateTime
	DeletedAt  gorm.DeletedAt

	Building DorBuilding `gorm:"foreignKey:Id;references:BuildingId"`
}

const (
	DorFloorIsPublicYes = 1
	DorFloorIsPublicNo  = 2
)
