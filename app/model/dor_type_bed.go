package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableDorTypeBed = "dor_type_bed"

type DorTypeBed struct {
	Id        int `gorm:"primary_key"`
	TypeId    int
	Name      string
	IsPublic  int8
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt
}

const (
	DorTypeBedIsPublicYes = 1
	DorTypeBedIsPublicNo  = 2
)
