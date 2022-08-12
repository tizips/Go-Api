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
	"saas/kernel/config"
	"saas/kernel/config/configs"
	"saas/kernel/data"
	"saas/kernel/snowflake"
	"strings"
	"time"
)

func DoUploadBySimple(ctx *gin.Context, dirs string, file *multipart.FileHeader) (upload *UploadBySimple, err error) {

	switch config.Values.File.Driver {
	case configs.FileDriverQiniu:
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

	if err := os.MkdirAll(config.Application.Runtime+filepath, 0750); err != nil {
		return nil, err
	}

	filename := snowflake.Snowflake.Generate().String() + path.Ext(file.Filename)

	filepath += "/" + filename

	if err := ctx.SaveUploadedFile(file, config.Application.Runtime+filepath); err != nil {
		return nil, err
	}

	return &UploadBySimple{
		Name: filename,
		Path: filepath,
		Url:  config.Values.Server.Url + filepath,
	}, nil
}

func doUploadBySimpleWithQiniu(ctx *gin.Context, dirs string, file *multipart.FileHeader) (*UploadBySimple, error) {

	redis := fmt.Sprintf("%s:qiniu:%s", config.Values.Server.Name, config.Values.File.QiniuAccess)

	token, _ := data.Redis.Get(ctx, redis).Result()

	if token == "" {

		policy := storage.PutPolicy{
			Scope:   config.Values.File.QiniuBucket,
			Expires: 7200,
		}

		mac := qbox.NewMac(config.Values.File.QiniuAccess, config.Values.File.QiniuSecret)

		token = policy.UploadToken(mac)

		if token != "" {
			data.Redis.Set(ctx, redis, token, time.Duration(policy.Expires)*time.Second)
		}

	}

	filename := snowflake.Snowflake.Generate().String() + path.Ext(file.Filename)

	key := dirs + "/" + filename

	if config.Values.File.QiniuPrefix != "" {
		key = "/" + config.Values.File.QiniuPrefix + key
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
		Url:  config.Values.File.QiniuDomain + "/" + key,
	}, nil
}

type UploadBySimple struct {
	Name string
	Path string
	Url  string
}
