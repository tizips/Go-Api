package basic

type ToBuildingByList struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Order     uint   `json:"order"`
	IsEnable  uint8  `json:"is_enable"`
	IsPublic  uint8  `json:"is_public"`
	CreatedAt string `json:"created_at"`
}

type ToBuildingByOnline struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	IsPublic uint8  `json:"is_public,omitempty"`
}
