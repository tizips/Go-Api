package auth

type ToRoleByPaginateResponse struct {
	Id          uint     `json:"id"`
	Name        string   `json:"name"`
	Summary     string   `json:"summary"`
	Permissions [][]uint `json:"permissions"`
	CreatedAt   string   `json:"created_at"`
}

type ToRoleByEnableResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
