package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableOaStaffLeave = "oa_staff_leave"

type OaStaffLeave struct {
	Id          uint `gorm:"primary_key"`
	MemberId    string
	StaffId     uint
	Channel     string
	OpenId      uint
	LastWorkDay carbon.DateTime
	ReasonType  string
	ReasonMemo  string
	PreStatus   string
	Status      string
	HandoverId  uint
	CreatedAt   carbon.DateTime
	UpdatedAt   carbon.DateTime
	DeletedAt   gorm.DeletedAt
}
