package users

type Service interface {
	Authenticate(string) error
	Create(string) error
	Delete(string) error
}

type userService struct{}

func NewService() (Service, error) {
	return &userService{}, nil
}

func (us *userService) Authenticate(userEmail string) error {

	return nil
}
func (us *userService) Create(userEmail string) error {

	return nil
}

func (us *userService) Delete(userEmail string) error {

	return nil
}
