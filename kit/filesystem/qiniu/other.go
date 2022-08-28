package qiniu

import (
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"saas/kernel/app"
	"time"
)

func (that *Qiniu) Mac() *qbox.Mac {

	return qbox.NewMac(app.Cfg.File.Qiniu.Access, app.Cfg.File.Qiniu.Secret)
}

func (that *Qiniu) Token() (token string) {

	token, _ = that.redis.Get(that.ctx, that.key).Result()

	if token == "" {

		policy := storage.PutPolicy{
			Scope:   app.Cfg.File.Qiniu.Bucket,
			Expires: 7200,
		}

		token = policy.UploadToken(that.Mac())

		if token != "" {
			that.redis.Set(that.ctx, that.key, token, time.Duration(policy.Expires)*time.Second)
		}

	}

	return token
}
