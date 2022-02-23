package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableOaStaff = "oa_staff"

type OaStaff struct {
	Id        uint `gorm:"primary_key"`
	MemberId  string
	No        string
	ManagerId uint
	Title     string
	Email     string
	HiredDate carbon.DateTime
	Remark    string
	Status    string
	IsEnable  uint8
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt
}
