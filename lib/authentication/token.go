package authentication

type AuthToken struct {
	T string `json:"token" form:"token"`
}
