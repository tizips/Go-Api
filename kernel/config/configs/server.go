package configs

type Server struct {
	Name string `default:"saas"`
	Mode string `default:"release"`
	Port int32  `default:"8080"`
	Url  string `default:"http://127.0.0.1:8080"`
	Node int64  `default:"0"`
	File string `default:"system"`
}

const (
	ServerFileSystem = "system"
	ServerFileQiniu  = "qiniu"
)
