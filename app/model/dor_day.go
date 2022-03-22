package model

import (
	"github.com/golang-module/carbon/v2"
)

const TableDorDay = "dor_day"

type DorDay struct {
	Id             uint `gorm:"primary_key"`
	CategoryId     uint
	TypeId         uint
	BuildingId     uint
	FloorId        uint
	RoomId         uint
	BedId          uint
	PeopleId       uint
	MemberId       string
	MasterPeopleId uint
	MasterMemberId string
	Date           carbon.Date
	CreatedAt      carbon.DateTime
}
