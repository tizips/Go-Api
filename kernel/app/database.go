package app

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	MySQL *gorm.DB
	Redis *redis.Client
)
