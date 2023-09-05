package login

import (
	"github.com/dgdraganov/fuzzy-user-api/pkg/model"
)

type Registry interface {
	LoginUser(dto model.LoginDTO) (string, error)
}
