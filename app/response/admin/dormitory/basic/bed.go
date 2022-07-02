package basic

type ToBedByPaginate struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Building  string `json:"building"`
	Floor     string `json:"floor"`
	Room      string `json:"room"`
	Order     int    `json:"order"`
	IsEnable  int8   `json:"is_enable"`
	IsPublic  int8   `json:"is_public"`
	CreatedAt string `json:"created_at"`
}

type ToBedByOnline struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	IsPublic int8   `json:"is_public,omitempty"`
}
