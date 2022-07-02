package stay

type ToPeopleByPaginate struct {
	Id            int      `json:"id"`
	Category      string   `json:"category"`
	Building      string   `json:"building"`
	Floor         string   `json:"floor"`
	Room          string   `json:"room"`
	Bed           string   `json:"bed"`
	Name          string   `json:"name"`
	Mobile        string   `json:"mobile"`
	Titles        string   `json:"titles,omitempty"`
	Staff         string   `json:"staff,omitempty"`
	Departments   []string `json:"departments,omitempty"`
	Manager       any      `json:"manager,omitempty"`
	Certification any      `json:"certification,omitempty"`
	IsTemp        int8     `json:"is_temp"`
	Start         string   `json:"start"`
	End           string   `json:"end,omitempty"`
	Remark        string   `json:"remark,omitempty"`
	CreatedAt     string   `json:"created_at"`
}

type ToPeopleByPaginateOfManager struct {
	Name   string `json:"name,omitempty"`
	Mobile string `json:"mobile,omitempty"`
}

type ToPeopleByPaginateOfCertification struct {
	No      string `json:"no,omitempty"`
	Address string `json:"address,omitempty"`
}
