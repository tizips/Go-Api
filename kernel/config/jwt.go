package config

type Jwt struct {
	Secret string
	Leeway string `default:"3"`
}
