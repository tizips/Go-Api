package app

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	Database  *gorm.DB
	Redis     *redis.Client
	Databases = make(map[string]*gorm.DB, 0)
)
