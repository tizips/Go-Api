package architecture

type ToModuleByList struct {
	Id        int    `json:"id"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	IsEnable  int8   `json:"is_enable"`
	Order     int    `json:"order"`
	CreatedAt string `json:"created_at"`
}

type ToModuleByOnline struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
