package system

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"os"
	"saas/kernel/config"
	"saas/kernel/data"
)

func Migrate() {
	fmt.Println("Migrate")

	m := connect()

	if err := m.Up(); err != nil {
		fmt.Printf("Mysql migrate fail:%s", err.Error())
		return
	}

	fmt.Println("Migrate success")
}

func connect() *migrate.Migrate {

	db, err := sql.Open(config.Values.Database.Driver, data.GetDns())
	if err != nil {
		fmt.Printf("Mysql connect fail:%s", err.Error())
		os.Exit(1)
	}

	defer db.Close()

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		fmt.Printf("Mysql casbin fail:%s", err.Error())
		os.Exit(1)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migration",
		"migrate",
		driver,
	)

	if err != nil {
		fmt.Printf("Mysql migrate fail:%s", err.Error())
		os.Exit(1)
	}

	return m
}
