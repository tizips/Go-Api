package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableDorGrant = "dor_grant"

type DorGrant struct {
	Id        int `gorm:"primary_key"`
	Object    string
	PackageId int
	Remark    string
	Cancel    string
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt

	Package     DorPackage       `gorm:"foreignKey:Id;references:PackageId"`
	BindDevices []DorGrantDevice `gorm:"foreignKey:GrantId;references:Id"`
}

// 	类型：package=打包；device=设备
const (
	DorGrantObjectPackage = "package"
	DorGrantObjectDevice  = "device"
)
