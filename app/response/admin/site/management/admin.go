package management

type ToAdminByPaginate struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Mobile   string `json:"mobile"`
	Roles    []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"roles"`
	IsEnable  int8   `json:"is_enable"`
	CreatedAt string `json:"created_at"`
}
