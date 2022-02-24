package stay

type ToCategoryByListResponse struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Order     uint   `json:"order"`
	IsTemp    uint8  `json:"is_temp"`
	IsEnable  uint8  `json:"is_enable"`
	CreatedAt string `json:"created_at"`
}

type ToCategoryByOnlineResponse struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	IsTemp uint8  `json:"is_temp"`
}
