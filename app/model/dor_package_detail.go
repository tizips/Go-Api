package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableDorPackageDetail = "dor_package_detail"

type DorPackageDetail struct {
	Id        uint `gorm:"primaryKey"`
	PackageId uint
	DeviceId  uint
	Number    uint
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt

	Device DorDevice `gorm:"foreignKey:Id;references:DeviceId"`
}
