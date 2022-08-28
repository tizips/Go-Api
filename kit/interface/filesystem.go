package _interface

import (
	"mime/multipart"
)

type FilesystemInterface interface {
	Upload() FilesystemInterface

	Mkdir(dir string) error

	Save(file *multipart.FileHeader, dir, name string) (uri, filename string, err error)

	List(dir string) ([]FilesystemFiles, error)

	Info(name string) error

	Delete(uri string) error

	Url(uri string) string

	Exist(uri string) bool
}

type FilesystemFiles struct {
	Name  string
	Url   string
	IsDir bool
}
