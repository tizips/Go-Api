package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableDorDevice = "dor_device"

type DorDevice struct {
	Id            int `gorm:"primary_key"`
	CategoryId    int
	No            string
	Name          string
	Specification string
	Price         int
	Unit          string
	Indemnity     int
	StockTotal    int
	StockUsed     int
	Remark        string
	CreatedAt     carbon.DateTime
	UpdatedAt     carbon.DateTime
	DeletedAt     gorm.DeletedAt

	Category DorAssetCategory `gorm:"foreignKey:Id;references:CategoryId"`
}
