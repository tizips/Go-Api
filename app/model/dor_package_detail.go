package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableDorPackageDetail = "dor_package_detail"

type DorPackageDetail struct {
	Id        int `gorm:"primaryKey"`
	PackageId int
	DeviceId  int
	Number    int
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt

	Device DorDevice `gorm:"foreignKey:Id;references:DeviceId"`
}
