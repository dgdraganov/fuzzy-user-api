package core_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dgdraganov/fuzzy-user-api/internal/core"
	"github.com/dgdraganov/fuzzy-user-api/pkg/model"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type jwtIssuerMock struct {
	generate func(*model.TokenInfo) *jwt.Token
	sign     func(token *jwt.Token) (string, error)
	validate func(token string) (jwt.MapClaims, error)
}

func (jwt *jwtIssuerMock) Generate(info *model.TokenInfo) *jwt.Token {
	return jwt.generate(info)
}
func (jwt *jwtIssuerMock) Sign(token *jwt.Token) (string, error) {
	return jwt.sign(token)
}
func (jwt *jwtIssuerMock) Validate(token string) (jwt.MapClaims, error) {
	return jwt.validate(token)
}

type repositoryMock struct {
	create  func(any) error
	getUser func(email string) (model.User, error)
}

func (r *repositoryMock) Create(entity any) error {
	return r.create(entity)
}
func (r *repositoryMock) GetUser(email string) (model.User, error) {
	return r.getUser(email)
}

func Test_UserExists_True(t *testing.T) {
	userEmail := "test@test.com"
	repo := &repositoryMock{
		getUser: func(email string) (model.User, error) {
			if email == userEmail {
				return model.User{}, nil
			}
			return model.User{}, errors.New("unexpected error")
		},
	}

	fuzzy := core.NewFuzzy(repo, nil)
	exists, err := fuzzy.UserExists(userEmail)

	var expectedError error = nil
	var expectedExists bool = true
	if err != nil {
		t.Fatalf("unexpected error, expected: %s, got: %s", expectedError, err)
	}

	if expectedExists != exists {
		t.Fatalf("failed user exists, expected: %t, got: %t", expectedExists, exists)
	}
}

func Test_UserExists_False(t *testing.T) {
	userEmail := "test@test.com"
	repo := &repositoryMock{
		getUser: func(email string) (model.User, error) {
			if email == userEmail {
				return model.User{}, fmt.Errorf("mock error: %w", gorm.ErrRecordNotFound)
			}
			return model.User{}, errors.New("unexpected error")
		},
	}

	fuzzy := core.NewFuzzy(repo, nil)
	exists, err := fuzzy.UserExists(userEmail)

	var expectedError error = nil
	var expectedExists bool = false
	if err != nil {
		t.Fatalf("unexpected error, expected: %s, got: %s", expectedError, err)
	}

	if expectedExists != exists {
		t.Fatalf("failed user exists, expected: %t, got: %t", expectedExists, exists)
	}
}

func Test_UserExists_Error(t *testing.T) {
	userEmail := "test@test.com"
	expectedError := errors.New("fake error")
	repo := &repositoryMock{
		getUser: func(email string) (model.User, error) {
			if email == userEmail {
				return model.User{}, expectedError
			}
			return model.User{}, errors.New("unexpected error")
		},
	}

	fuzzy := core.NewFuzzy(repo, nil)
	exists, err := fuzzy.UserExists(userEmail)

	var expectedExists bool = false
	if !errors.Is(err, expectedError) {
		t.Fatalf("unexpected error, expected: %s, got: %s", expectedError, err)
	}

	if expectedExists != exists {
		t.Fatalf("failed user exists, expected: %t, got: %t", expectedExists, exists)
	}
}

// func Test_LoginUser_Success(t *testing.T) {

// 	dto := model.LoginDTO{
// 		Email:    "test@test.com",
// 		Password: "testPass",
// 	}

// 	expectedToken := "fake_token"

// 	repoMock := &repositoryMock{
// 		getUser: func(email string) (model.User, error) {
// 			if email == dto.Email {
// 				return model.User{
// 					FirstName:    "Tanio",
// 					LastName:     "Tanev",
// 					Email:        dto.Email,
// 					PasswordHash: "fake_pass_hash",
// 				}, nil
// 			}
// 			return model.User{}, errors.New("unexpected error")
// 		},
// 	}

// 	jwtMock := &jwtIssuerMock{
// 		generate: func(ti *model.TokenInfo) *jwt.Token {

// 			return nil
// 		},
// 		sign: func(token *jwt.Token) (string, error) {
// 			bcrypt.
// 			return expectedToken, nil
// 		},
// 	}

// 	fuzzy := core.NewFuzzy(repoMock, jwtMock)

// 	token, err := fuzzy.LoginUser(dto)
// 	if err != nil {
// 		t.Fatalf("unexpected error: %s", err)
// 	}

// 	if token != expectedToken {
// 		t.Fatalf("token does not match, expected: %s, got: %s", expectedToken, token)
// 	}
// }
