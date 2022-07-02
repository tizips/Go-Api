package model

import (
	"github.com/golang-module/carbon/v2"
)

const TableDorDay = "dor_day"

type DorDay struct {
	Id             int `gorm:"primary_key"`
	CategoryId     int
	TypeId         int
	BuildingId     int
	FloorId        int
	RoomId         int
	BedId          int
	PeopleId       int
	MemberId       string
	MasterPeopleId int
	MasterMemberId string
	Date           carbon.Date
	CreatedAt      carbon.DateTime
}
