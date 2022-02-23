package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableOaStaffOpen = "oa_staff_open"

type OaStaffOpen struct {
	Id        uint `gorm:"primary_key"`
	MemberId  string
	StaffId   uint
	Channel   string
	Openid    string
	Unionid   string
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt
}
