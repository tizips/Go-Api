package response

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type Responses struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []any  `json:"data"`
}

type Paginate[T any] struct {
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
	Data  []T   `json:"data"`
}
