package basic

type Enable struct {
	IsEnable uint8 `form:"is_enable" json:"is_enable" binding:"eq=1|eq=2"`
}
