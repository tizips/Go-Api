package asset

type ToGrantByPaginateResponse struct {
	Id        uint                                `json:"id"`
	Package   string                              `json:"package,omitempty"`
	Devices   []ToGrantByPaginateOfDeviceResponse `json:"devices,omitempty"`
	Remark    string                              `json:"remark"`
	CreatedAt string                              `json:"created_at"`
}

type ToGrantByPaginateOfDeviceResponse struct {
	Name   string `json:"name"`
	Number uint   `json:"number"`
}
