package basic

type ToBuildingByList struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
	IsEnable  int8   `json:"is_enable"`
	IsPublic  int8   `json:"is_public"`
	CreatedAt string `json:"created_at"`
}

type ToBuildingByOnline struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	IsPublic int8   `json:"is_public,omitempty"`
}
