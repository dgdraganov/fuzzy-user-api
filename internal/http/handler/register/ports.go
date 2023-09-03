package register

type UserRepository interface {
	Create(obj any) error
}
