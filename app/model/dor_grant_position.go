package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableDorGrantPosition = "dor_grant_device"

type DorGrantPosition struct {
	Id         int `gorm:"primary_key"`
	GrantId    int
	Object     string
	TypeId     int
	TypeBedId  int
	BuildingId int
	FloorId    int
	RoomId     int
	BedId      int
	CreatedAt  carbon.DateTime
	UpdatedAt  carbon.DateTime
	DeletedAt  gorm.DeletedAt
}

const (
	DorGrantPositionLive     = "live"
	DorGrantPositionType     = "type"
	DorGrantPositionBuilding = "building"
	DorGrantPositionFloor    = "floor"
	DorGrantPositionRoom     = "room"
	DorGrantPositionBed      = "bed"
)
