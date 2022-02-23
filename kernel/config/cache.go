package config

type Cache struct {
	Prefix string `default:"cache"`
	Ttl    string `default:"86400"`
}
