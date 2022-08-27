package cmd

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"os"
	"saas/kernel/database"
)

func Migrate(cmd *cobra.Command) error {

	return nil
}

func handler() {
	fmt.Println("Migrate")

	m := connect()

	if err := m.Up(); err != nil {
		fmt.Printf("MySQL migrate fail:%s", err.Error())
		return
	}

	fmt.Println("Migrate success")
}

func connect() *migrate.Migrate {

	db, err := sql.Open("mysql", database.MySQLDail())
	if err != nil {
		fmt.Printf("MySQL connect fail:%s", err.Error())
		os.Exit(1)
	}

	defer db.Close()

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		fmt.Printf("MySQL casbin fail:%s", err.Error())
		os.Exit(1)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migration",
		"migrate",
		driver,
	)

	if err != nil {
		fmt.Printf("MySQL migrate fail:%s", err.Error())
		os.Exit(1)
	}

	return m
}
