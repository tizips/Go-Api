package model

type SysCasbin struct {
	Id    uint   `gorm:"primary_key"`
	PType string `gorm:"column:ptype"`
	V0    string
	V1    string
	V2    string
	V3    string
	V4    string
	V5    string
}
