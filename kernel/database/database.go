package database

import (
	"errors"
	"fmt"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	gl "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"os"
	"saas/kernel/app"
	"saas/kernel/logger"
	"time"
)

func InitDatabase() {

	var err error

	app.Database, err = NewDatabase(app.Cfg.Database.Driver)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func Connect(driver string) (data *gorm.DB, err error) {

	var ok = false

	if data, ok = app.Databases[driver]; ok {
		return data, nil
	}

	if data, err = NewDatabase(driver); err != nil {
		return
	}

	return data, nil
}

func NewDatabase(driver string) (data *gorm.DB, err error) {

	log := logger.NewGormLogger()
	log.SlowThreshold = time.Millisecond * 200

	var prefix = ""
	var dialectal gorm.Dialector

	switch driver {
	case "mysql":
		dialectal = mysql.Open(md())
		prefix = app.Cfg.Database.MySQL.Prefix
	case "postgresql":
		dialectal = postgres.Open(pd())
		prefix = app.Cfg.Database.PostgreSQL.Prefix
	case "sqlserver":
		dialectal = sqlserver.Open(sd())
		prefix = app.Cfg.Database.SQLServer.Prefix
	case "clickhouse":
		dialectal = clickhouse.Open(cd())
		prefix = app.Cfg.Database.Clickhouse.Prefix
	default:
		return nil, errors.New("driver fail")
	}

	data, err = gorm.Open(dialectal, &gorm.Config{
		Logger: log.LogMode(gl.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		return
	}

	db, err := data.DB()

	if err != nil {
		return
	}

	// SetMaxIdleCons 设置空闲连接池中连接的最大数量
	// SetMaxOpenCons 设置打开数据库连接的最大数量。
	// SetConnMaxLifetime 设置了连接可复用的最大时间。

	switch driver {
	case "mysql":
		db.SetMaxIdleConns(app.Cfg.Database.MySQL.MaxIdle)
		db.SetMaxOpenConns(app.Cfg.Database.MySQL.MaxIdle)
		db.SetConnMaxLifetime(time.Duration(app.Cfg.Database.MySQL.MaxLifetime) * time.Second)
	case "postgresql":
		db.SetMaxIdleConns(app.Cfg.Database.PostgreSQL.MaxIdle)
		db.SetMaxOpenConns(app.Cfg.Database.PostgreSQL.MaxIdle)
		db.SetConnMaxLifetime(time.Duration(app.Cfg.Database.PostgreSQL.MaxLifetime) * time.Second)
	case "sqlserver":
		db.SetMaxIdleConns(app.Cfg.Database.SQLServer.MaxIdle)
		db.SetMaxOpenConns(app.Cfg.Database.SQLServer.MaxIdle)
		db.SetConnMaxLifetime(time.Duration(app.Cfg.Database.SQLServer.MaxLifetime) * time.Second)
	case "clickhouse":
		db.SetMaxIdleConns(app.Cfg.Database.Clickhouse.MaxIdle)
		db.SetMaxOpenConns(app.Cfg.Database.Clickhouse.MaxIdle)
		db.SetConnMaxLifetime(time.Duration(app.Cfg.Database.Clickhouse.MaxLifetime) * time.Second)
	}

	return data, nil
}

// MySQL
func md() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
		app.Cfg.Database.MySQL.Username,
		app.Cfg.Database.MySQL.Password,
		app.Cfg.Database.MySQL.Host,
		app.Cfg.Database.MySQL.Port,
		app.Cfg.Database.MySQL.Database,
		app.Cfg.Database.MySQL.Charset)
}

// PostgreSQL
func pd() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		app.Cfg.Database.PostgreSQL.Host,
		app.Cfg.Database.PostgreSQL.Username,
		app.Cfg.Database.PostgreSQL.Password,
		app.Cfg.Database.PostgreSQL.Database,
		app.Cfg.Database.PostgreSQL.Port,
		app.Cfg.Database.PostgreSQL.SslMode,
		app.Cfg.Database.PostgreSQL.Timezone,
	)
}

// SQL Server
func sd() string {
	return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		app.Cfg.Database.SQLServer.Username,
		app.Cfg.Database.SQLServer.Password,
		app.Cfg.Database.SQLServer.Host,
		app.Cfg.Database.SQLServer.Port,
		app.Cfg.Database.SQLServer.Database,
	)
}

// Clickhouse
func cd() string {
	return fmt.Sprintf("tcp://%s:%d?database=%s&username=%s&password=%s&read_timeout=%d&write_timeout=%d",
		app.Cfg.Database.Clickhouse.Host,
		app.Cfg.Database.Clickhouse.Port,
		app.Cfg.Database.Clickhouse.Database,
		app.Cfg.Database.Clickhouse.Username,
		app.Cfg.Database.Clickhouse.Password,
		app.Cfg.Database.Clickhouse.ReadTimeout,
		app.Cfg.Database.Clickhouse.WriteTimeout,
	)
}
