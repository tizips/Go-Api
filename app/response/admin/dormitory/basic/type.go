package basic

type ToTypeByList struct {
	Id        uint                `json:"id"`
	Name      string              `json:"name"`
	Beds      []ToTypeByListOfBed `json:"beds,omitempty"`
	Order     uint                `json:"order"`
	IsEnable  uint8               `json:"is_enable"`
	CreatedAt string              `json:"created_at"`
}

type ToTypeByListOfBed struct {
	Name     string `json:"name"`
	IsPublic uint8  `json:"is_public"`
}

type ToTypeByOnline struct {
	Id   uint                  `json:"id"`
	Name string                `json:"name"`
	Beds []ToTypeByOnlineOfBed `json:"beds,omitempty"`
}

type ToTypeByOnlineOfBed struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
