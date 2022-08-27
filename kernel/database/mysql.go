package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"os"
	"saas/kernel/app"
	"saas/kernel/logger"
	"time"
)

func InitMySQL() {

	log := logger.NewGormLogger()
	log.SlowThreshold = time.Millisecond * 200

	var err error

	app.MySQL, err = gorm.Open(mysql.Open(MySQLDail()), &gorm.Config{
		Logger: log.LogMode(gormLogger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   app.Cfg.Database.MySQL.Prefix,
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		fmt.Printf("MySQL load fail:%s", err.Error())
		os.Exit(1)
	}

	sqlDB, err := app.MySQL.DB()

	if err != nil {
		fmt.Printf("MySQL load fail:%s", err.Error())
		os.Exit(1)
	}

	// SetMaxIdleCons 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(app.Cfg.Database.MySQL.MaxIdle)

	// SetMaxOpenCons 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(app.Cfg.Database.MySQL.MaxIdle)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Duration(app.Cfg.Database.MySQL.MaxLifetime) * time.Second)
}

func MySQLDail() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		app.Cfg.Database.MySQL.Username,
		app.Cfg.Database.MySQL.Password,
		app.Cfg.Database.MySQL.Host,
		app.Cfg.Database.MySQL.Port,
		app.Cfg.Database.MySQL.Database,
		app.Cfg.Database.MySQL.Charset)
}
