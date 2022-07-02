package basic

type ToRoomByPaginate struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Building  string `json:"building"`
	Floor     string `json:"floor"`
	Type      string `json:"type"`
	TypeId    int    `json:"type_id"`
	Order     int    `json:"order"`
	IsFurnish int8   `json:"is_furnish"`
	IsEnable  int8   `json:"is_enable"`
	IsPublic  int8   `json:"is_public"`
	CreatedAt string `json:"created_at"`
}

type ToRoomByOnline struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	IsPublic int8   `json:"is_public,omitempty"`
}
