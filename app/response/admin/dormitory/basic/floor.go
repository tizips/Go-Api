package basic

type ToFloorByList struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Building  string `json:"building"`
	Order     int    `json:"order"`
	IsEnable  int8   `json:"is_enable"`
	IsPublic  int8   `json:"is_public"`
	CreatedAt string `json:"created_at"`
}

type ToFloorByOnline struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	IsPublic int8   `json:"is_public"`
}
