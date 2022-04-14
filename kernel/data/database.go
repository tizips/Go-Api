package data

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"os"
	"saas/kernel/config"
	"saas/kernel/logger"
	"time"
)

var Database *gorm.DB

func InitDatabase() {

	log := logger.NewGormLogger()
	log.SlowThreshold = time.Millisecond * 200

	var err error

	Database, err = gorm.Open(mysql.Open(GetDns()), &gorm.Config{
		Logger: log.LogMode(gormLogger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config.Values.Database.Prefix,
			SingularTable: true,
		},
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
	sqlDB.SetMaxIdleConns(config.Values.Database.MaxIdle)

	// SetMaxOpenCons 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(config.Values.Database.MaxIdle)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Duration(config.Values.Database.MaxLifetime) * time.Second)
}

func GetDns() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.Values.Database.Username,
		config.Values.Database.Password,
		config.Values.Database.Host,
		config.Values.Database.Port,
		config.Values.Database.Database,
		config.Values.Database.Charset)
}
