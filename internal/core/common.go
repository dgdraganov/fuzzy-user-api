package core

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type fuzzy struct {
	repo      Repository
	jwtIssuer JwtIssuer
}

// NewFuzzy is a constructor function for the fuzzy type
func NewFuzzy(db Repository, issuer JwtIssuer) *fuzzy {
	return &fuzzy{
		repo:      db,
		jwtIssuer: issuer,
	}
}

func (f *fuzzy) UserExists(email string) (bool, error) {
	_, err := f.repo.GetUser(email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("repo get user: %w", err)
	}
	return true, nil
}
