package crontab

import (
	"github.com/robfig/cron/v3"
	"saas/app/crontab/admin/dormitory"
)

var crontab *cron.Cron

func InitCrontab() {

	crontab = cron.New()

	register()

	crontab.Start()

}

func register() {

	dormitory.CrontabDayPeople(crontab)

}
