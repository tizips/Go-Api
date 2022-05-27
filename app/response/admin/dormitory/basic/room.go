package basic

type ToRoomByPaginate struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Building  string `json:"building"`
	Floor     string `json:"floor"`
	Type      string `json:"type"`
	TypeId    uint   `json:"type_id"`
	Order     uint   `json:"order"`
	IsFurnish uint8  `json:"is_furnish"`
	IsEnable  uint8  `json:"is_enable"`
	IsPublic  uint8  `json:"is_public"`
	CreatedAt string `json:"created_at"`
}

type ToRoomByOnline struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	IsPublic uint8  `json:"is_public,omitempty"`
}
