package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableDorRoom = "dor_room"

type DorRoom struct {
	Id         int `gorm:"primary_key"`
	BuildingId int
	FloorId    int
	TypeId     int
	Name       string
	Order      int
	IsFurnish  int8
	IsEnable   int8
	IsPublic   int8
	CreatedAt  carbon.DateTime
	UpdatedAt  carbon.DateTime
	DeletedAt  gorm.DeletedAt

	Building DorBuilding `gorm:"foreignKey:Id;references:BuildingId"`
	Floor    DorFloor    `gorm:"foreignKey:Id;references:FloorId"`
	Type     DorType     `gorm:"foreignKey:Id;references:TypeId"`
}

const (
	DorRoomIsFurnishYes = 1
	DorRoomIsFurnishNo  = 2

	DorRoomIsPublicYes = 1
	DorRoomIsPublicNo  = 2
)
