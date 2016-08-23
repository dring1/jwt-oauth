package models

type Model interface {
	Crud()
	Read()
	Update()
	Delete()
}
