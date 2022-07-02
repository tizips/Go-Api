package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableDorStayCategory = "dor_stay_category"

type DorStayCategory struct {
	Id        int `gorm:"primary_key"`
	Name      string
	Order     int
	IsTemp    int8
	IsEnable  int8
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt
}

const (
	DorStayCategoryIsTempYes = 1
	DorStayCategoryIsTempNo  = 0
)
