package configs

type Cache struct {
	Prefix string `default:"cache"`
	Ttl    int32  `default:"86400"`
}
