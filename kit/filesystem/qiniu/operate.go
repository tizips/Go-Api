package qiniu

import (
	"bufio"
	"github.com/gookit/goutil/dump"
	"github.com/qiniu/go-sdk/v7/storage"
	"mime/multipart"
	"path/filepath"
	"saas/kernel/app"
	_interface "saas/kit/interface"
	"strings"
)

func (that *Qiniu) Mkdir(dir string) error {
	return nil
}

func (that *Qiniu) Save(file *multipart.FileHeader, dir, name string) (uri, filename string, err error) {

	var fo multipart.File

	if fo, err = file.Open(); err != nil {
		return
	}

	defer fo.Close()

	if strings.HasPrefix(name, "/") {
		name = strings.Trim(name, "/")
	}

	if name == "" {
		name = app.Snowflake.Generate().String() + filepath.Ext(file.Filename)
	}

	uri = dir + "/" + name

	uri = that.path(uri)

	resume := storage.NewFormUploader(nil)

	var ret storage.PutRet

	err = resume.Put(that.ctx, &ret, that.Token(), uri, bufio.NewReader(fo), file.Size, nil)

	if err != nil {
		return
	}

	return uri, name, err
}

func (that *Qiniu) List(path string) ([]_interface.FilesystemFiles, error) {

	path = that.path(path)

	dump.P(path)

	files := make([]_interface.FilesystemFiles, 0)

	dirs := make(map[string]string, 0)
	children := make([]_interface.FilesystemFiles, 0)

	manager := storage.NewBucketManager(that.Mac(), nil)

	marker := ""

	for {

		entries, _, s, b, err := manager.ListFiles(app.Cfg.File.Qiniu.Bucket, path, "", marker, 100)

		if err != nil {
			return nil, err
		}

		for _, item := range entries {

			names := strings.Split(strings.TrimLeft(strings.TrimPrefix(item.Key, path), "/"), "/")

			if len(names) > 1 {
				dirs[names[0]] = names[0]
			} else if len(names) == 1 {
				children = append(children, _interface.FilesystemFiles{Name: names[0], Url: that.Url(item.Key), IsDir: false})
			}
		}

		if b {
			marker = s
		} else {
			break
		}
	}

	for _, item := range dirs {
		files = append(files, _interface.FilesystemFiles{Name: item, IsDir: true})
	}

	files = append(files, children...)

	return files, nil
}

func (that *Qiniu) Delete(uri string) error {

	manager := storage.NewBucketManager(that.Mac(), nil)

	return manager.Delete(app.Cfg.File.Qiniu.Bucket, uri)
}

func (that *Qiniu) Exist(uri string) bool {
	return true
}

func (that *Qiniu) Url(key string) string {
	return storage.MakePublicURLv2(app.Cfg.File.Qiniu.Domain, key)
}

func (that *Qiniu) path(path string) string {

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if app.Cfg.File.Qiniu.Prefix != "" {
		path = app.Cfg.File.Qiniu.Prefix + path
	}

	path = strings.Trim(path, "/")

	return path

}
