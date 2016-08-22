package model

type Model interface {
	Crud()
	Read()
	Update()
	Delete()
}
