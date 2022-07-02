package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableDorGrantDetail = "dor_grant_detail"

type DorGrantDetail struct {
	Id         int `gorm:"primary_key"`
	GrantId    int
	PackageId  int
	PositionId int
	TypeId     int
	BuildingId int
	FloorId    int
	RoomId     int
	BedId      int
	PeopleId   int
	MemberId   string `gorm:"default:null"`
	DeviceId   int
	Number     int
	IsPublic   int8
	CreatedAt  carbon.DateTime
	UpdatedAt  carbon.DateTime
	DeletedAt  gorm.DeletedAt
}
