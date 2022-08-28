package local

import (
	"saas/kit/interface"
)

type Local struct {
	_interface.FilesystemInterface

	prefix string
}

func New() _interface.FilesystemInterface {
	return new(Local)
}

func (that *Local) Upload() _interface.FilesystemInterface {

	that.prefix = "upload"

	return that
}
