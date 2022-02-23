package basic

type ToFloorByListResponse struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Building  string `json:"building"`
	Order     uint   `json:"order"`
	IsEnable  uint8  `json:"is_enable"`
	CreatedAt string `json:"created_at"`
}

type ToFloorByOnlineResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
