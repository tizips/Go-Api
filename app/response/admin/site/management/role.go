package management

type ToRoleByPaginate struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Summary     string  `json:"summary"`
	Permissions [][]int `json:"permissions"`
	CreatedAt   string  `json:"created_at"`
}

type ToRoleByEnable struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
