package auth

type ToRoleByPaginate struct {
	Id          uint     `json:"id"`
	Name        string   `json:"name"`
	Summary     string   `json:"summary"`
	Permissions [][]uint `json:"permissions"`
	CreatedAt   string   `json:"created_at"`
}

type ToRoleByEnable struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
