package configs

type File struct {
	Driver      string `default:"local"`
	QiniuAccess string
	QiniuSecret string
	QiniuBucket string
	QiniuDomain string
	QiniuPrefix string
}

const (
	FileDriverQiniu = "qiniu"
)
