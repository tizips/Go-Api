package basic

type DoLoginByAccess struct {
	Token    string `json:"token"`
	ExpireAt int64  `json:"expire_at"`
}
