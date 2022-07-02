package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableMemCertificationImage = "mem_certification_image"

type MemCertificationImage struct {
	Id              int `gorm:"primary_key"`
	MemberId        string
	CertificationId int
	Url             string
	CreatedAt       carbon.DateTime
	UpdatedAt       carbon.DateTime
	DeletedAt       gorm.DeletedAt
}
