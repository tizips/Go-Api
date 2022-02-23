package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
	"saas/kernel/cache"
)

const TableSysAdmin = "sys_admin"

type SysAdmin struct {
	Id        uint `gorm:"primary_key"`
	Username  string
	Mobile    string
	Email     string
	Nickname  string
	Avatar    string
	Password  string
	IsEnable  uint8
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt

	BindRoles []SysAdminBindRole `gorm:"foreignKey:AdminId"`

	cache.Model
}
