package basic

type Paginate struct {
	Page uint `form:"page" json:"page" binding:"omitempty,number,gte=0"`
	Size uint `form:"size" json:"size" binding:"omitempty,number,gte=0,lte=20"`
}

func (p *Paginate) GetPage() uint {
	if p.Page > 0 {
		return p.Page
	}
	return 1
}

func (p *Paginate) GetSize() uint {
	if p.Size > 0 {
		return p.Size
	}
	return 15
}

func (p *Paginate) GetLimit() int {
	return int(p.GetSize())
}

func (p *Paginate) GetOffset() int {
	return int((p.GetPage() - 1) * p.GetSize())
}
