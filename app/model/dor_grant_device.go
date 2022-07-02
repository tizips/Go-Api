package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableDorGrantDevice = "dor_grant_device"

type DorGrantDevice struct {
	Id        int `gorm:"primary_key"`
	GrantId   int
	DeviceId  int
	Number    int
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt

	Device DorDevice `gorm:"foreignKey:Id;references:DeviceId"`
}
