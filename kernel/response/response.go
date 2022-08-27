package response

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"database"`
}

type Responses struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []any  `json:"database"`
}

type Paginate struct {
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
	Data  []any `json:"database"`
}
