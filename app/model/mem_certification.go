package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableMemCertification = "mem_certification"

type MemCertification struct {
	Id         int `gorm:"primary_key"`
	MemberId   string
	Status     string
	Type       string
	Name       string
	No         string
	Other      string
	ValidStart carbon.DateTime
	ValidEnd   carbon.DateTime
	CreatedAt  carbon.DateTime
	UpdatedAt  carbon.DateTime
	DeletedAt  gorm.DeletedAt

	Images []MemCertificationImage `gorm:"foreignKey:MemberId"`
}
