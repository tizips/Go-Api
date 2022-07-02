package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableDorPackage = "dor_package"

type DorPackage struct {
	Id        int `gorm:"primaryKey"`
	Name      string
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt

	Details []DorPackageDetail `gorm:"foreignKey:PackageId;references:Id"`
}
