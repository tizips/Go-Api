package asset

type ToPackageByPaginate struct {
	Id        int                            `json:"id"`
	Name      string                         `json:"name"`
	Devices   []ToPackageByPaginateOfDevices `json:"devices"`
	CreatedAt string                         `json:"created_at"`
}

type ToPackageByPaginateOfDevices struct {
	Category int    `json:"category"`
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Number   int    `json:"number"`
}

type ToPackageByOnline struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
