package asset

type ToDeviceByPaginate struct {
	Id            int    `json:"id"`
	Category      string `json:"category"`
	CategoryId    int    `json:"category_id"`
	Name          string `json:"name"`
	No            string `json:"no"`
	Specification string `json:"specification"`
	Price         int    `json:"price"`
	Unit          string `json:"unit"`
	Indemnity     int    `json:"indemnity"`
	StockTotal    int    `json:"stock_total"`
	StockUsed     int    `json:"stock_used"`
	Remark        string `json:"remark"`
	CreatedAt     string `json:"created_at"`
}

type ToDeviceByOnline struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
