package manage

type ToRoleByPaginate struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Summary     string  `json:"summary"`
	Permissions [][]int `json:"permissions"`
	CreatedAt   string  `json:"created_at"`
}

type ToRoleByOnline struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
