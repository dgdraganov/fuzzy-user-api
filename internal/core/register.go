package core

import (
	"fmt"

	"github.com/dgdraganov/fuzzy-user-api/pkg/model"
	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

func (f *fuzzy) RegisterUser(dto model.RegisterDTO) error {

	err := validateRegisterDTO(dto)
	if err != nil {
		return fmt.Errorf("validate register dto: %w", err)
	}

	user, err := f.prepareUserRegister(dto)
	if err != nil {
		return fmt.Errorf("prapare user register: %w", err)
	}

	if err := f.repo.Create(&user); err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (f *fuzzy) prepareUserRegister(dto model.RegisterDTO) (model.User, error) {
	var res model.User
	hs, err := bcrypt.GenerateFromPassword([]byte(dto.Password), 10)
	if err != nil {
		return model.User{}, fmt.Errorf("bcrypt generate password hash: %w", err)
	}
	res.FirstName = dto.FirstName
	res.LastName = dto.LastName
	res.Email = dto.Email
	res.PasswordHash = string(hs)
	return res, nil
}

func validateRegisterDTO(dto model.RegisterDTO) error {
	validate := validator.New()
	err := validate.Struct(dto)
	if err != nil {
		return fmt.Errorf("validate struct: %w", err)
	}
	return nil
}
