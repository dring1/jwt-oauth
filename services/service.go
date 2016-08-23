package services

type Service interface {
	RegisterService(*[]Service)
}
