package qiniu

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"saas/kernel/app"
	_interface "saas/kit/interface"
)

type Qiniu struct {
	_interface.FilesystemInterface

	ctx   context.Context
	key   string
	redis *redis.Client
}

func New() *Qiniu {

	qn := &Qiniu{
		ctx:   context.Background(),
		key:   fmt.Sprintf("%s:token:qiniu:%s", app.Cfg.Server.Name, app.Cfg.File.Qiniu.Access),
		redis: app.Redis,
	}

	return qn
}

func (that *Qiniu) Upload() _interface.FilesystemInterface {

	return that
}

func (that *Qiniu) SetContext(ctx context.Context) {
	that.ctx = ctx
}

func (that *Qiniu) SetKey(key string) error {

	if key == "" {
		return errors.New("key not exist")
	}

	that.key = key

	return nil
}

func (that *Qiniu) SetRedis(rds *redis.Client) {
	that.redis = rds
}
