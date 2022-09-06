package cmd

import (
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gookit/color"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"os"
	"saas/kernel/app"
)

type migrate struct {
	//	父级命令
	command *cobra.Command

	//	子级命令
	cmd *cobra.Command

	path string
}

func NewMigrate(command *cobra.Command) {

	m := migrate{
		command: command,
		path:    "migration",
	}

	m.init()

}

func (that *migrate) init() {

	that.cmd = &cobra.Command{
		Use:   "migrate",
		Short: "数据迁移",
	}

	that.goose()

	that.create()
	that.make()
	that.rollback()
	that.redo()
	that.status()
	that.version()

	that.command.AddCommand(that.cmd)

}

func (that *migrate) goose() {

	goose.SetTableName(app.Cfg.Database.MySQL.Prefix + "sys_migration")

	if err := goose.SetDialect(app.Cfg.Database.Driver); err != nil {
		color.Errorln(err)
		os.Exit(0)
	}

}

func (that *migrate) create() {

	cmd := &cobra.Command{
		Use:   "create",
		Short: "创建迁移文件",
		Run: func(c *cobra.Command, args []string) {

			name, _ := c.Flags().GetString("name")

			db, _ := app.Database.DB()

			err := goose.Create(db, that.path, name, "sql")

			if err != nil {
				color.Errorf("create migration fail: %v", err)
			}
		},
	}

	cmd.Flags().StringP("name", "n", "", "表名")
	_ = cmd.MarkFlagRequired("name")

	that.cmd.AddCommand(cmd)
}

func (that *migrate) make() {

	cmd := &cobra.Command{
		Use:   "make",
		Short: "执行迁移",
		Long:  "默认将数据库迁移到可用的最新版本",
		Run: func(c *cobra.Command, args []string) {

			db, _ := app.Database.DB()

			var err error

			if one, _ := c.Flags().GetBool("one"); one {
				err = goose.UpByOne(db, that.path)
			} else if version, _ := c.Flags().GetInt64("version"); version > 0 {
				err = goose.UpTo(db, that.path, version)
			} else {
				err = goose.Up(db, that.path)
			}

			if err != nil {
				color.Errorln(err)
			}

		},
	}

	cmd.Flags().BoolP("one", "o", false, "执行一次迁移文件")
	cmd.Flags().Int64P("version", "v", 0, "将数据库迁移到特定版本")

	that.cmd.AddCommand(cmd)
}

func (that *migrate) rollback() {

	cmd := &cobra.Command{
		Use:   "rollback",
		Short: "迁移回撤",
		Run: func(c *cobra.Command, args []string) {

			db, _ := app.Database.DB()

			var err error

			if version, _ := c.Flags().GetInt64("version"); version > 0 {
				err = goose.DownTo(db, that.path, version)
			} else if all, _ := c.Flags().GetBool("all"); all {
				err = goose.Reset(db, that.path)
			} else {
				err = goose.Down(db, that.path)
			}

			if err != nil {
				color.Errorln(err)
			}

		},
	}

	cmd.Flags().BoolP("all", "a", false, "回滚所有迁移")
	cmd.Flags().Int64P("version", "v", 0, "将数据库回撤到特定版本")

	that.cmd.AddCommand(cmd)
}

func (that *migrate) redo() {

	cmd := &cobra.Command{
		Use:   "redo",
		Short: "重新运行最新的迁移",
		Run: func(c *cobra.Command, args []string) {

			db, _ := app.Database.DB()

			if err := goose.Redo(db, that.path); err != nil {
				color.Errorln(err)
			}
		},
	}

	that.cmd.AddCommand(cmd)
}

func (that *migrate) status() {

	cmd := &cobra.Command{
		Use:   "status",
		Short: "查看迁移文件",
		Run: func(c *cobra.Command, args []string) {

			db, _ := app.Database.DB()

			if err := goose.Status(db, that.path); err != nil {
				color.Errorln(err)
			}
		},
	}

	that.cmd.AddCommand(cmd)
}

func (that *migrate) version() {

	cmd := &cobra.Command{
		Use:   "version",
		Short: "查看迁移版本",
		Run: func(c *cobra.Command, args []string) {

			db, _ := app.Database.DB()

			err := goose.Version(db, that.path)

			if err != nil {
				color.Errorln(err)
			}
		},
	}

	that.cmd.AddCommand(cmd)
}
