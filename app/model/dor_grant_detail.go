package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableDorGrantDetail = "dor_grant_detail"

type DorGrantDetail struct {
	Id         uint `gorm:"primary_key"`
	GrantId    uint
	PackageId  uint
	PositionId uint
	TypeId     uint
	BuildingId uint
	FloorId    uint
	RoomId     uint
	BedId      uint
	PeopleId   uint
	MemberId   string `gorm:"default:null"`
	DeviceId   uint
	Number     uint
	IsPublic   uint8
	CreatedAt  carbon.DateTime
	UpdatedAt  carbon.DateTime
	DeletedAt  gorm.DeletedAt
}
