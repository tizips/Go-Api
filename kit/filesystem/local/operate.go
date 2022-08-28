package local

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"saas/kernel/app"
	_interface "saas/kit/interface"
	"strings"
)

func (that *Local) Mkdir(dir string) error {

	var ok bool

	if ok = that.Exist(dir); !ok {
		if err := os.MkdirAll(app.Dir.Runtime+that.path(dir), 0750); err != nil {
			return err
		}
	}

	return nil
}

func (that *Local) Save(file *multipart.FileHeader, dir, name string) (uri, filename string, err error) {

	if file == nil {
		return "", "", errors.New("文件读取失败")
	}

	if err = that.Mkdir(dir); err != nil {
		return "", "", err
	}

	if strings.HasPrefix(name, "/") {
		name = strings.Trim(name, "/")
	}

	if name == "" {
		name = app.Snowflake.Generate().String() + filepath.Ext(file.Filename)
	}

	fp := that.path(dir) + "/" + name

	var src multipart.File

	if src, err = file.Open(); err != nil {
		return "", "", err
	}

	defer src.Close()

	var out *os.File

	if out, err = os.Create(app.Dir.Runtime + fp); err != nil {
		return "", "", err
	}

	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return "", "", err
	}

	return fp, name, nil
}

func (that *Local) List(dir string) ([]_interface.FilesystemFiles, error) {

	var ok bool

	if ok = that.Exist(dir); !ok {
		return nil, errors.New("file not exist")
	}

	dir = that.path(dir)

	dirs, err := os.ReadDir(app.Dir.Runtime + dir)

	if err != nil {
		return nil, err
	}

	files := make([]_interface.FilesystemFiles, 0)

	for _, item := range dirs {

		file := _interface.FilesystemFiles{Name: item.Name(), IsDir: item.IsDir()}

		if !item.IsDir() {
			file.Url = that.Url(dir + "/" + item.Name())
		}

		files = append(files, file)
	}

	return files, nil
}

func (that *Local) Delete(uri string) error {

	var ok bool

	if ok = that.Exist(uri); !ok {
		return errors.New("file not exist")
	}

	return os.Remove(app.Dir.Runtime + that.path(uri))
}

func (that *Local) Url(uri string) string {

	return app.Cfg.Server.Url + uri
}

func (that *Local) Exist(uri string) bool {

	if _, err := os.Stat(app.Dir.Runtime + that.path(uri)); os.IsNotExist(err) {
		return false
	}

	return true
}

func (that *Local) path(path string) string {

	if path != "" && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if that.prefix != "" {
		path = "/" + strings.TrimLeft(that.prefix, "/") + path
	}

	return path
}
