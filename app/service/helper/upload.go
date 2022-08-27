package helper

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"mime/multipart"
	"os"
	"path"
	"saas/kernel/app"
	"strings"
	"time"
)

func DoUploadBySimple(ctx *gin.Context, dirs string, file *multipart.FileHeader) (upload *UploadBySimple, err error) {

	switch app.Cfg.File.Driver {
	case app.Cfg.File.Driver:
		upload, err = doUploadBySimpleWithQiniu(ctx, dirs, file)
	default:
		upload, err = doUploadBySimpleWithSystem(ctx, dirs, file)
	}

	return upload, err
}

func doUploadBySimpleWithSystem(ctx *gin.Context, dirs string, file *multipart.FileHeader) (*UploadBySimple, error) {

	filepath := "/upload"

	if !strings.HasPrefix(dirs, "/") {
		filepath += "/"
	}

	filepath += dirs

	if err := os.MkdirAll(app.Dir.Runtime+filepath, 0750); err != nil {
		return nil, err
	}

	filename := app.Snowflake.Generate().String() + path.Ext(file.Filename)

	filepath += "/" + filename

	if err := ctx.SaveUploadedFile(file, app.Dir.Runtime+filepath); err != nil {
		return nil, err
	}

	return &UploadBySimple{
		Name: filename,
		Path: filepath,
		Url:  app.Cfg.Server.Url + filepath,
	}, nil
}

func doUploadBySimpleWithQiniu(ctx *gin.Context, dirs string, file *multipart.FileHeader) (*UploadBySimple, error) {

	redis := fmt.Sprintf("%s:qiniu:%s", app.Cfg.Server.Name, app.Cfg.File.Qiniu.Access)

	token, _ := app.Redis.Get(ctx, redis).Result()

	if token == "" {

		policy := storage.PutPolicy{
			Scope:   app.Cfg.File.Qiniu.Bucket,
			Expires: 7200,
		}

		mac := qbox.NewMac(app.Cfg.File.Qiniu.Access, app.Cfg.File.Qiniu.Secret)

		token = policy.UploadToken(mac)

		if token != "" {
			app.Redis.Set(ctx, redis, token, time.Duration(policy.Expires)*time.Second)
		}

	}

	filename := app.Snowflake.Generate().String() + path.Ext(file.Filename)

	key := dirs + "/" + filename

	if app.Cfg.File.Qiniu.Prefix != "" {
		key = "/" + app.Cfg.File.Qiniu.Prefix + key
	}

	if strings.HasPrefix(key, "/") {
		key = string([]rune(key)[1:])
	}

	resume := storage.NewFormUploader(nil)

	f, err := file.Open()
	if err != nil {
		return nil, err
	}

	var ret storage.PutRet

	err = resume.Put(ctx, &ret, token, key, bufio.NewReader(f), file.Size, nil)
	if err != nil {
		return nil, err
	}

	return &UploadBySimple{
		Name: filename,
		Path: "/" + key,
		Url:  app.Cfg.File.Qiniu.Domain + "/" + key,
	}, nil
}

type UploadBySimple struct {
	Name string
	Path string
	Url  string
}
