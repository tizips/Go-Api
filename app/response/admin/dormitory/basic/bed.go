package basic

type ToBedByPaginateResponse struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Building  string `json:"building"`
	Floor     string `json:"floor"`
	Room      string `json:"room"`
	Order     uint   `json:"order"`
	IsEnable  uint8  `json:"is_enable"`
	IsPublic  uint8  `json:"is_public"`
	CreatedAt string `json:"created_at"`
}

type ToBedByOnlineResponse struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	IsPublic uint8  `json:"is_public,omitempty"`
}
