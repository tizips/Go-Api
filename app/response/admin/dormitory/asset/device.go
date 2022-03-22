package asset

type ToDeviceByPaginateResponse struct {
	Id            uint   `json:"id"`
	Category      string `json:"category"`
	CategoryId    uint   `json:"category_id"`
	Name          string `json:"name"`
	No            string `json:"no"`
	Specification string `json:"specification"`
	Price         uint   `json:"price"`
	Unit          string `json:"unit"`
	Indemnity     uint   `json:"indemnity"`
	StockTotal    uint   `json:"stock_total"`
	StockUsed     uint   `json:"stock_used"`
	Remark        string `json:"remark"`
	CreatedAt     string `json:"created_at"`
}

type ToDeviceByOnlineResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
