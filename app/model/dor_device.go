package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableDorDevice = "dor_device"

type DorDevice struct {
	Id            uint `gorm:"primary_key"`
	CategoryId    uint
	No            string
	Name          string
	Specification string
	Price         uint
	Unit          string
	Indemnity     uint
	StockTotal    uint
	StockUsed     uint
	Remark        string
	CreatedAt     carbon.DateTime
	UpdatedAt     carbon.DateTime
	DeletedAt     gorm.DeletedAt

	Category DorAssetCategory `gorm:"foreignKey:Id;references:CategoryId"`
}
