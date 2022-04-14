package configs

type Jwt struct {
	Secret   string
	Leeway   int64 `default:"3"`
	Lifetime int   `default:"12"`
}
