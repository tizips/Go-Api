package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
	"saas/kernel/cache"
)

const TableSysAdmin = "sys_admin"

type SysAdmin struct {
	Id        int    `gorm:"primary_key"`
	Username  string `gorm:"default:null"`
	Mobile    string `gorm:"default:null"`
	Email     string `gorm:"default:null"`
	Nickname  string
	Avatar    string
	Password  string
	IsEnable  int8
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt

	BindRoles []SysAdminBindRole `gorm:"foreignKey:AdminId"`

	cache.Model
}
