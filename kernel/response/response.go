package response

type Response struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type Responses struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
	Data    []any  `json:"data"`
}
