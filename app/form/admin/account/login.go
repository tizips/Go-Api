package account

type DoLoginForm struct {
	Username string `form:"username" json:"username" binding:"required,min=4,max=20"`
	Password string `form:"password" json:"password" binding:"required,min=4,max=20"`
}
