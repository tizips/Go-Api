package asset

type ToGrantByPaginate struct {
	Id        int                         `json:"id"`
	Package   string                      `json:"package,omitempty"`
	Devices   []ToGrantByPaginateOfDevice `json:"devices,omitempty"`
	Remark    string                      `json:"remark"`
	CreatedAt string                      `json:"created_at"`
}

type ToGrantByPaginateOfDevice struct {
	Name   string `json:"name"`
	Number int    `json:"number"`
}
