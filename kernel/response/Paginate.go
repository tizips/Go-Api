package response

type Paginate struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Total int64 `json:"total"`
		Page  uint  `json:"page"`
		Size  uint  `json:"size"`
		Data  []any `json:"data"`
	} `json:"data"`
}
