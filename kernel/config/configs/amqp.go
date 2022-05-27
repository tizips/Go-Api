package configs

type Amqp struct {
	Host      string `default:"127.0.0.1"`
	Port      int32  `default:"5672"`
	Username  string `default:"admin"`
	Password  string
	Vhost     string `default:"/"`
	Reconnect int32  `default:"3"`
}
