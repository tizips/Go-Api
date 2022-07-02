package configs

type Redis struct {
	Host        string `default:"127.0.0.1"`
	Auth        string
	Port        int32  `default:"6379"`
	Db          int    `default:"0"`
	CachePrefix string `default:"cache"`
	CacheTtl    int32  `default:"86400"`
}
