package auth

type ToAdminByPaginateResponse struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Mobile   string `json:"mobile"`
	Roles    []struct {
		Id   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"roles"`
	IsEnable  uint8  `json:"is_enable"`
	CreatedAt string `json:"created_at"`
}
