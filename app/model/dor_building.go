package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableDorBuilding = "dor_building"

type DorBuilding struct {
	Id        int `gorm:"primary_key"`
	Name      string
	Order     int
	IsEnable  int8
	IsPublic  int8
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt
}

const (
	DorBuildingIsPublicYes = 1
	DorBuildingIsPublicNo  = 2
)
