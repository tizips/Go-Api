package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableDorRoom = "dor_room"

type DorRoom struct {
	Id         uint `gorm:"primary_key"`
	BuildingId uint
	FloorId    uint
	TypeId     uint
	Name       string
	Order      uint
	IsFurnish  uint8
	IsEnable   uint8
	CreatedAt  carbon.DateTime
	UpdatedAt  carbon.DateTime
	DeletedAt  gorm.DeletedAt

	Building DorBuilding `gorm:"foreignKey:Id;references:BuildingId"`
	Floor    DorFloor    `gorm:"foreignKey:Id;references:FloorId"`
	Type     DorType     `gorm:"foreignKey:Id;references:TypeId"`
}
