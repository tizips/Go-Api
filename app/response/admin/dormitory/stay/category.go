package stay

type ToCategoryByList struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
	IsTemp    int8   `json:"is_temp"`
	IsEnable  int8   `json:"is_enable"`
	CreatedAt string `json:"created_at"`
}

type ToCategoryByOnline struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	IsTemp int8   `json:"is_temp"`
}
