package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableMemGroup = "mem_group"

type MemGroup struct {
	Id        uint `gorm:"primary_key"`
	Code      string
	Name      string
	IsDefault uint8
	IsEnable  uint8
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt
}
