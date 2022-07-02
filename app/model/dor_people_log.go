package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableDorPeopleLog = "dor_people_log"

type DorPeopleLog struct {
	Id        int `gorm:"primary_key"`
	PeopleId  int
	MemberId  string
	Status    string
	Detail    string
	Remark    string
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt
}

const (
	DorPeopleLogStatusLive     = "live"
	DorPeopleLogStatusLeave    = "leave"
	DorPeopleLogStatusChange   = "change"
	DorPeopleLogStatusRefill   = "refill"
	DorPeopleLogStatusPositive = "positive"
)
