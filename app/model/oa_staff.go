package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableOaStaff = "oa_staff"

type OaStaff struct {
	Id        int `gorm:"primary_key"`
	MemberId  string
	No        string `gorm:"default:null"`
	ManagerId int
	Title     string
	Email     string `gorm:"default:null"`
	HiredDate carbon.DateTime
	Remark    string
	Status    string
	IsEnable  int8
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt
}
