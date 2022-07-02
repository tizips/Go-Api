package model

import (
	"github.com/golang-module/carbon/v2"
)

const TableDorPeople = "dor_people"

type DorPeople struct {
	Id         int `gorm:"primary_key"`
	CategoryId int
	BuildingId int
	FloorId    int
	RoomId     int
	BedId      int
	TypeId     int
	MemberId   string
	MasterId   int
	Start      carbon.Date  `gorm:"default:null"`
	End        *carbon.Date `gorm:"default:null"`
	Remark     string
	IsTemp     int8
	Status     string
	CreatedAt  carbon.DateTime
	UpdatedAt  carbon.DateTime

	Member        *MemMember        `gorm:"foreignKey:Id;references:MemberId"`
	Staff         *OaStaff          `gorm:"foreignKey:MemberId;references:MemberId"`
	Certification *MemCertification `gorm:"foreignKey:MemberId;references:MemberId"`
	Category      *DorStayCategory  `gorm:"foreignKey:Id;references:CategoryId"`
	Building      *DorBuilding      `gorm:"foreignKey:Id;references:BuildingId"`
	Floor         *DorFloor         `gorm:"foreignKey:Id;references:FloorId"`
	Room          *DorRoom          `gorm:"foreignKey:Id;references:RoomId"`
	Bed           *DorBed           `gorm:"foreignKey:Id;references:BedId"`
	Type          *DorType          `gorm:"foreignKey:Id;references:TypeId"`
	Master        *DorPeople        `gorm:"foreignKey:MasterId;references:Id"`
}

const (
	DorPeopleStatusLive  = "live"
	DorPeopleStatusLeave = "leave"

	DorPeopleIsTempYes = 1
	DorPeopleIsTempNo  = 2

	DorPeopleIsMasterYes = 1
	DorPeopleIsMasterNo  = 2
)

func (m *DorPeople) GetStatusName() string {
	name := "未知"
	switch m.Status {
	case DorPeopleStatusLive:
		name = "在住"
	case DorPeopleStatusLeave:
		name = "离宿"
	}
	return name
}
