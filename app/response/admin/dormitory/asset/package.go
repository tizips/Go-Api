package asset

type ToPackageByPaginateResponse struct {
	Id        uint                                   `json:"id"`
	Name      string                                 `json:"name"`
	Devices   []ToPackageByPaginateOfDevicesResponse `json:"devices"`
	CreatedAt string                                 `json:"created_at"`
}

type ToPackageByPaginateOfDevicesResponse struct {
	Category uint   `json:"category"`
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Number   uint   `json:"number"`
}

type ToPackageByOnlineResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
