package config

type Redis struct {
	Host string `default:"127.0.0.1"`
	Auth string `default:"omitempty"`
	Port string `default:"6379"`
	Db   string `default:"0"`
}
