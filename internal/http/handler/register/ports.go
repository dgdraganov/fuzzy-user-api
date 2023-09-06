package register

import "github.com/dgdraganov/fuzzy-user-api/pkg/model"

type Registry interface {
	RegisterUser(model.RegisterDTO) error
	UserExists(email string) (bool, error)
}
