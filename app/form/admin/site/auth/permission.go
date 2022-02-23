package auth

type ToPermissionByTreeForm struct {
	Module uint `form:"module" binding:"required,number,gt=0"`
}

type DoPermissionByCreateForm struct {
	Module uint   `form:"module" json:"module" binding:"required,gt=0"`
	Parent uint   `form:"parent" json:"parent" binding:"gte=0"`
	Name   string `form:"name" json:"name" binding:"required,min=2,max=20"`
	Slug   string `form:"slug" json:"slug" binding:"required,min=2,max=64"`
	Method string `form:"method" json:"method" binding:"omitempty,required_with=Path,oneof=GET POST PUT DELETE"`
	Path   string `form:"path" json:"path" binding:"omitempty,required_with=Method,max=64"`
}

type DoPermissionByUpdateForm struct {
	Module uint   `form:"module" json:"module" binding:"required,gt=0"`
	Parent uint   `form:"parent" json:"parent" binding:"gte=0"`
	Name   string `form:"name" json:"name" binding:"required,min=2,max=20"`
	Slug   string `form:"slug" json:"slug" binding:"required,min=2,max=64"`
	Method string `form:"method" json:"method" binding:"omitempty,required_with=Path,oneof=GET POST PUT DELETE"`
	Path   string `form:"path" json:"path" binding:"omitempty,required_with=Method,max=64"`
}
