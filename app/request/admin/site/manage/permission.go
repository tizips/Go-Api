package manage

type ToPermissionByTree struct {
	Module int `form:"module" binding:"required,number,gt=0" label:"模块"`
}

type DoPermissionByCreate struct {
	Module int    `form:"module" json:"module" binding:"required,gt=0" label:"模块"`
	Parent int    `form:"parent" json:"parent" binding:"gte=0" label:"父级"`
	Name   string `form:"name" json:"name" binding:"required,min=2,max=20" label:"名称"`
	Slug   string `form:"slug" json:"slug" binding:"required,min=2,max=64" label:"标识"`
	Method string `form:"method" json:"method" binding:"omitempty,required_with=Path,oneof=GET POST PUT DELETE"`
	Path   string `form:"path" json:"path" binding:"omitempty,required_with=Method,max=64"`
}

type DoPermissionByUpdate struct {
	Module int    `form:"module" json:"module" binding:"required,gt=0" label:"模块"`
	Parent int    `form:"parent" json:"parent" binding:"gte=0" label:"父级"`
	Name   string `form:"name" json:"name" binding:"required,min=2,max=20" label:"名称"`
	Slug   string `form:"slug" json:"slug" binding:"required,min=2,max=64" label:"标识"`
	Method string `form:"method" json:"method" binding:"omitempty,required_with=Path,oneof=GET POST PUT DELETE"`
	Path   string `form:"path" json:"path" binding:"omitempty,required_with=Method,max=64"`
}
