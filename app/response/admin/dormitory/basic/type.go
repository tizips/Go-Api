package basic

type ToTypeByList struct {
	Id        int                 `json:"id"`
	Name      string              `json:"name"`
	Beds      []ToTypeByListOfBed `json:"beds,omitempty"`
	Order     int                 `json:"order"`
	IsEnable  int8                `json:"is_enable"`
	CreatedAt string              `json:"created_at"`
}

type ToTypeByListOfBed struct {
	Name     string `json:"name"`
	IsPublic int8   `json:"is_public"`
}

type ToTypeByOnline struct {
	Id   int                   `json:"id"`
	Name string                `json:"name"`
	Beds []ToTypeByOnlineOfBed `json:"beds,omitempty"`
}

type ToTypeByOnlineOfBed struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
