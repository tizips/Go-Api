package configs

type Database struct {
	Driver      string `default:"mysql"`
	Host        string `default:"127.0.0.1"`
	Port        int32  `default:"3306"`
	Database    string `default:"upper"`
	Username    string `default:"root"`
	Password    string `default:"omitempty"`
	Charset     string `default:"utf8mb4"`
	Prefix      string `default:""`
	MaxIdle     int    `default:"10"`
	MaxOpen     int    `default:"100"`
	MaxLifetime int    `default:"60"`
}
