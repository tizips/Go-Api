package helper

type ToApiByList struct {
	Module uint `form:"module" json:"module" binding:"required,gt=0"`
}
