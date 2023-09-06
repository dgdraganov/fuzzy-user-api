package core

import (
	"errors"
	"fmt"

	"github.com/dgdraganov/fuzzy-user-api/pkg/model"
	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidPassword error = errors.New("invalid password")

func (f *fuzzy) LoginUser(dto model.LoginDTO) (string, error) {
	err := validateLoginDTO(dto)
	if err != nil {
		return "", fmt.Errorf("validate login dto: %w", err)
	}

	user, err := f.repo.GetUser(dto.Email)
	if err != nil {
		return "", fmt.Errorf("get password hash: %w", err)
	}

	info, err := f.prepareTokenInfo(dto, user)
	if err != nil {
		return "", fmt.Errorf("prapare token info: %w", err)
	}

	token := f.jwtIssuer.Generate(&info)
	tokenStr, err := f.jwtIssuer.Sign(token)
	if err != nil {
		return "", fmt.Errorf("token signing: %w", err)
	}
	return tokenStr, nil
}

func validateLoginDTO(dto model.LoginDTO) error {
	validate := validator.New()
	err := validate.Struct(dto)
	if err != nil {
		return fmt.Errorf("validate struct: %w", err)
	}
	return nil
}

func (f *fuzzy) prepareTokenInfo(dto model.LoginDTO, user model.User) (model.TokenInfo, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(dto.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return model.TokenInfo{}, ErrInvalidPassword
	}
	if err != nil {
		return model.TokenInfo{}, fmt.Errorf("compare pass and pass hash: %w", err)
	}

	var res model.TokenInfo
	res.Email = user.Email
	res.FirstName = user.FirstName
	res.LastName = user.LastName
	res.Subject = "Login"
	res.Expiration = 24

	return res, nil
}

type TokenInfo struct {
	Email      string
	UserID     int
	FirstName  string
	LastName   string
	Subject    string
	Expiration int
}
