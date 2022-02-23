package config

type Database struct {
	Driver      string `default:"mysql"`
	Host        string `default:"127.0.0.1"`
	Port        string `default:"3306"`
	Database    string `default:"upper"`
	Username    string `default:"root"`
	Password    string `default:"omitempty"`
	Charset     string `default:"utf8mb4"`
	Prefix      string `default:""`
	MaxIdle     string `default:"10"`
	MaxOpen     string `default:"100"`
	MaxLifetime string `default:"60"`
}
