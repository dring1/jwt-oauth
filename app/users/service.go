package users

type Service interface {
	Authenticate(string) error
	Create()
	Delete()
}
