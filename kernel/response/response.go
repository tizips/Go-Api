package response

type Response struct {
	Code    uint        `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Responses struct {
	Code    uint          `json:"code"`
	Message string        `json:"message"`
	Data    []interface{} `json:"data"`
}
