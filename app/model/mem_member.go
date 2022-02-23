package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
	"saas/kernel/cache"
)

const TableMemMember = "mem_member"

type MemMember struct {
	Id        string `gorm:"primary_key"`
	GroupId   uint
	Username  string
	Mobile    string
	Email     string
	Name      string
	Avatar    string
	Nickname  string
	Password  string
	Sex       uint8
	Province  uint
	City      uint
	Area      uint
	Year      uint16
	Month     uint8
	Day       uint8
	IsEnable  uint8
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt

	Group         MemGroup         `gorm:"references:GroupId"`
	Certification MemCertification `gorm:"foreignKey:CertificationId"`

	cache.Model
}
