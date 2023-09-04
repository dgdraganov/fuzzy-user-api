package core

type Repository interface {
	Create(any) error
}
