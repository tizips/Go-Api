package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
	"saas/kernel/cache"
)

const TableMemMember = "mem_member"

type MemMember struct {
	Id        string `gorm:"primary_key"`
	GroupId   int    `gorm:"default:0"`
	Username  string `gorm:"default:null"`
	Mobile    string `gorm:"default:null"`
	Email     string `gorm:"default:null"`
	Name      string
	Avatar    string
	Nickname  string
	Password  string
	Sex       int8  `gorm:"default:0"`
	Province  int   `gorm:"default:0"`
	City      int   `gorm:"default:0"`
	Area      int   `gorm:"default:0"`
	Year      int16 `gorm:"default:0"`
	Month     int8  `gorm:"default:0"`
	Day       int8  `gorm:"default:0"`
	IsEnable  int8  `gorm:"default:0"`
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt

	Group         MemGroup         `gorm:"references:GroupId;foreignKey:Id"`
	Certification MemCertification `gorm:"references:Id;foreignKey:MemberId"`

	cache.Model
}
