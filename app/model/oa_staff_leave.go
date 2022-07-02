package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableOaStaffLeave = "oa_staff_leave"

type OaStaffLeave struct {
	Id          int `gorm:"primary_key"`
	MemberId    string
	StaffId     int
	Channel     string
	OpenId      int
	LastWorkDay carbon.DateTime
	ReasonType  string
	ReasonMemo  string
	PreStatus   string
	Status      string
	HandoverId  int
	CreatedAt   carbon.DateTime
	UpdatedAt   carbon.DateTime
	DeletedAt   gorm.DeletedAt
}
