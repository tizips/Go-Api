package architecture

type ToModuleByList struct {
	Id        uint   `json:"id"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	IsEnable  uint8  `json:"is_enable"`
	Order     uint   `json:"order"`
	CreatedAt string `json:"created_at"`
}

type ToModuleByOnline struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
