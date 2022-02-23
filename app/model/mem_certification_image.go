package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableMemCertificationImage = "mem_certification_image"

type MemCertificationImage struct {
	Id              uint `gorm:"primary_key"`
	MemberId        string
	CertificationId uint
	Url             string
	CreatedAt       carbon.DateTime
	UpdatedAt       carbon.DateTime
	DeletedAt       gorm.DeletedAt
}
