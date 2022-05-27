package asset

type ToCategoryByList struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Order     uint   `json:"order"`
	IsEnable  uint8  `json:"is_enable"`
	CreatedAt string `json:"created_at"`
}

type ToCategoryByOnline struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
