package basic

type ToAccountByPermission struct {
	Module int `form:"module" binding:"required,numeric,gt=0" label:"模块"`
}

type DoAccountByUpdate struct {
	Avatar   string `json:"avatar" form:"avatar" binding:"required,url,max=255" label:"头像"`
	Password string `json:"password" form:"password" binding:"omitempty,min=6,max=20,password" label:"密码"`
}
