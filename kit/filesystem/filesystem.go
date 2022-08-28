package filesystem

import (
	"saas/kernel/app"
	"saas/kit/filesystem/local"
	"saas/kit/filesystem/qiniu"
	"saas/kit/interface"
)

func New() _interface.FilesystemInterface {
	return Disk(app.Cfg.File.Driver)
}

func Disk(disk string) _interface.FilesystemInterface {

	var storage _interface.FilesystemInterface

	switch disk {
	case "qiniu":
		storage = qiniu.New()
	default:
		storage = local.New()
	}

	return storage
}
