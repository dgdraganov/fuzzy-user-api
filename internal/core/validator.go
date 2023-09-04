package core

import "github.com/dgdraganov/fuzzy-user-api/pkg/model"

type RegisterValidator interface {
	Validate(dto model.RegisterDTO) error
}
