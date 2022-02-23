package basic

type Enable struct {
	IsEnable uint8 `form:"is_enable" json:"is_enable" binding:"oneof=0 1"`
}
