package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableDorGrantPosition = "dor_grant_device"

type DorGrantPosition struct {
	Id         uint `gorm:"primary_key"`
	GrantId    uint
	Object     string
	TypeId     uint
	TypeBedId  uint
	BuildingId uint
	FloorId    uint
	RoomId     uint
	BedId      uint
	CreatedAt  carbon.DateTime
	UpdatedAt  carbon.DateTime
	DeletedAt  gorm.DeletedAt
}

const (
	DorGrantPositionType     = "type"
	DorGrantPositionBuilding = "building"
	DorGrantPositionFloor    = "floor"
	DorGrantPositionRoom     = "room"
	DorGrantPositionBed      = "bed"
)
