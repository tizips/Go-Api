package local

import (
	_interface "saas/kit/interface"
)

type Local struct {
	_interface.FilesystemInterface

	prefix string
}

func New() *Local {
	return new(Local)
}

func (that *Local) Upload() _interface.FilesystemInterface {

	that.prefix = "upload"

	return that
}
