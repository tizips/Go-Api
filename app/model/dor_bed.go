package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableDorBed = "dor_bed"

type DorBed struct {
	Id         int `gorm:"primary_key"`
	BuildingId int
	FloorId    int
	RoomId     int
	TypeId     int
	BedId      int
	Name       string
	Order      int
	IsEnable   int8
	IsPublic   int8
	CreatedAt  carbon.DateTime
	UpdatedAt  carbon.DateTime
	DeletedAt  gorm.DeletedAt

	Building DorBuilding `gorm:"foreignKey:Id;references:BuildingId"`
	Floor    DorFloor    `gorm:"foreignKey:Id;references:FloorId"`
	Room     DorRoom     `gorm:"foreignKey:Id;references:RoomId"`
	Type     DorType     `gorm:"foreignKey:Id;references:TypeId"`
}

const (
	DorBedIsPublicYes = 1
	DorBedIsPublicNo  = 2
)
