package basic

type ToTypeByListResponse struct {
	Id        uint                        `json:"id"`
	Name      string                      `json:"name"`
	Beds      []ToTypeByListOfBedResponse `json:"beds,omitempty"`
	Order     uint                        `json:"order"`
	IsEnable  uint8                       `json:"is_enable"`
	CreatedAt string                      `json:"created_at"`
}

type ToTypeByListOfBedResponse struct {
	Name     string `json:"name"`
	IsPublic uint8  `json:"is_public"`
}

type ToTypeByOnlineResponse struct {
	Id   uint                          `json:"id"`
	Name string                        `json:"name"`
	Beds []ToTypeByOnlineOfBedResponse `json:"beds,omitempty"`
}

type ToTypeByOnlineOfBedResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
