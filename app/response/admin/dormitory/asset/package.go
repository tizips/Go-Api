package asset

type ToPackageByPaginate struct {
	Id        uint                           `json:"id"`
	Name      string                         `json:"name"`
	Devices   []ToPackageByPaginateOfDevices `json:"devices"`
	CreatedAt string                         `json:"created_at"`
}

type ToPackageByPaginateOfDevices struct {
	Category uint   `json:"category"`
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Number   uint   `json:"number"`
}

type ToPackageByOnline struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
