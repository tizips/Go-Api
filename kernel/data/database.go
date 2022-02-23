package data

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"os"
	"saas/kernel/config"
	"strconv"
	"time"
)

var Database *gorm.DB

func InitDatabase() {

	var err error

	Database, err = gorm.Open(mysql.Open(GetDns()), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		NamingStrategy:                           schema.NamingStrategy{TablePrefix: config.Configs.Database.Prefix, SingularTable: true},
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		fmt.Printf("Mysql load fail:%s", err.Error())
		os.Exit(1)
	}

	sqlDB, err := Database.DB()

	if err != nil {
		fmt.Printf("Mysql load fail:%s", err.Error())
		os.Exit(1)
	}

	// SetMaxIdleCons 设置空闲连接池中连接的最大数量
	maxIdle, _ := strconv.Atoi(config.Configs.Database.MaxIdle)
	sqlDB.SetMaxIdleConns(maxIdle)

	// SetMaxOpenCons 设置打开数据库连接的最大数量。
	maxOpen, _ := strconv.Atoi(config.Configs.Database.MaxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	maxLifetime, _ := strconv.Atoi(config.Configs.Database.MaxLifetime)
	sqlDB.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)
}

func GetDns() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		config.Configs.Database.Username,
		config.Configs.Database.Password,
		config.Configs.Database.Host,
		config.Configs.Database.Port,
		config.Configs.Database.Database,
		config.Configs.Database.Charset)
}
