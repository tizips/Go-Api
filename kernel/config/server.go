package config

type Server struct {
	Name string `default:"saas"`
	Mode string `default:"release"`
	Port string `default:"8080"`
	Url  string `default:"http://127.0.0.1:8080"`
}
