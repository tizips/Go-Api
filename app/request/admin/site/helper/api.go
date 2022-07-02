package helper

type ToApiByList struct {
	Module int `form:"module" json:"module" binding:"required,gt=0" label:"模块"`
}
