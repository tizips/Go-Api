package basic

type ToAccountByPermission struct {
	Module uint `form:"module" binding:"required,numeric,gt=0"`
}

type DoAccountByUpdate struct {
	Avatar   string `json:"avatar" form:"avatar" binding:"required,url,max=255"`
	Password string `json:"password" form:"password" binding:"omitempty,min=6,max=20,password"`
}
