package app

var Cfg struct {
	Server struct {
		Name string `yaml:"name" default:"saas"`
		Mode string `yaml:"mode" default:"release"`
		Port int32  `yaml:"port" default:"8080"`
		Url  string `yaml:"url" default:"http://127.0.0.1:8080"`
		Node int64  `yaml:"node" default:"1"`
	} `yaml:"server"`

	Database struct {
		MySQL struct {
			Host        string `yaml:"host" default:"127.0.0.1"`
			Port        int32  `yaml:"port" default:"3306"`
			Database    string `yaml:"database" default:"saas"`
			Username    string `yaml:"username" default:"root"`
			Password    string `yaml:"password"`
			Charset     string `yaml:"charset"`
			Collation   string `yaml:"collation"`
			Prefix      string `yaml:"prefix"`
			MaxIdle     int    `yaml:"max_idle" default:"10"`
			MaxOpen     int    `yaml:"max_open" default:"100"`
			MaxLifetime int    `yaml:"max_lifetime" default:"60"`
		} `yaml:"mysql"`
		Redis struct {
			Host        string `yaml:"host" default:"127.0.0.1"`
			Port        int32  `yaml:"port" default:"6379"`
			Password    string `yaml:"password"`
			Db          int    `yaml:"db" default:"0"`
			MaxConnAge  int    `yaml:"max_conn_age" default:"100"`
			PoolTimeout int    `yaml:"pool_timeout" default:"3"`
			IdleTimeout int    `yaml:"idle_timeout" default:"60"`
		} `yaml:"redis"`
	} `yaml:"database"`

	File struct {
		Driver string `yaml:"driver" default:"local"`
		Qiniu  struct {
			Access string `yaml:"access"`
			Secret string `yaml:"secret"`
			Bucket string `yaml:"bucket"`
			Domain string `yaml:"domain"`
			Prefix string `yaml:"prefix"`
		} `yaml:"qiniu"`
	} `yaml:"file"`

	Queue struct {
		Driver string `yaml:"driver" default:"amqp"`
		Amqp   struct {
			Host     string `yaml:"host" default:"127.0.0.1"`
			Port     int32  `yaml:"port" default:"5672"`
			Username string `yaml:"username" default:"admin"`
			Password string `yaml:"password"`
			Vhost    string `yaml:"vhost" default:"/"`
		} `yaml:"amqp"`
	} `yaml:"queue"`

	Jwt struct {
		Secret   string `yaml:"secret"`
		Leeway   int64  `yaml:"leeway" default:"3"`
		Lifetime int    `yaml:"lifetime" default:"12"`
	} `yaml:"jwt"`
}
