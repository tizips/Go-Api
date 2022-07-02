package basic

type DoUploadBySimple struct {
	Dir string `form:"dir" json:"dir" binding:"required,dir,max=20" label:"文件夹"`
}
