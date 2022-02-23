package helper

type ToApiByListForm struct {
	Module uint `form:"module" json:"module" binding:"required,gt=0"`
}
