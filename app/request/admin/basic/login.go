package basic

type DoLoginByAccess struct {
	Username string `form:"username" json:"username" binding:"required,username" label:"用户名"`
	Password string `form:"password" json:"password" binding:"required,password" label:"密码"`
}
