package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableDorBed = "dor_bed"

type DorBed struct {
	Id         uint `gorm:"primary_key"`
	BuildingId uint
	FloorId    uint
	RoomId     uint
	TypeId     uint
	BedId      uint
	Name       string
	Order      uint
	IsEnable   uint8
	IsPublic   uint8
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
