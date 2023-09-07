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
func Test_RegisterUser_Success(t *testing.T) {
	dto := model.RegisterDTO{
		FirstName: "Penko",
		LastName:  "Penkov",
		Email:     "test@gmail.com",
		Password:  "testPass",
	}
	var expected error

	repoMock := repositoryMock{
		create: func(a any) error {
			return nil
		},
	}
	fuzzy := core.NewFuzzy(&repoMock, nil)
	err := fuzzy.RegisterUser(dto)

	if err != expected {
		t.Fatalf("unexpected error, expected: %s, got: %s", expected, err)
	}
}
func Test_RegisterUser_Failed(t *testing.T) {
	dto := model.RegisterDTO{
		FirstName: "Penko",
		LastName:  "Penkov",
		Email:     "test@gmainl.com",
		Password:  "testPass",
	}

	var expected error = errors.New("expected mock error")
	repoMock := repositoryMock{
		create: func(a any) error {
			us, ok := a.(*model.User)
			fmt.Println(us)
			if ok && us.Email == dto.Email {
				return expected
			}
			return errors.New("unexpected error")
		},
	}

	fuzzy := core.NewFuzzy(&repoMock, nil)
	err := fuzzy.RegisterUser(dto)

	if !errors.Is(err, expected) {
		t.Fatalf("unexpected error, expected: %s, got: %s", expected, err)
	}
}

//	func (f *fuzzy) VerifyUser(jwtToken string) (map[string]any, error) {
//		claims, err := f.jwtIssuer.Validate(jwtToken)
//		if err != nil { return nil, fmt.Errorf("jwt validate: %w", err) }
//		return map[string]any(claims), nil
//	}
func Test_VerifyUser_Success(t *testing.T) {
	tokenMock := "fake_token"
	expectedClaims := map[string]any{
		"claim1": "value1",
		"claim2": 2,
	}
	jwtMock := jwtIssuerMock{
		validate: func(token string) (jwt.MapClaims, error) {
			if token == tokenMock {
				return expectedClaims, nil
			}
			return nil, errors.New("unexpected error")
		},
	}

	fuzzyMock := core.NewFuzzy(nil, &jwtMock)

	claims, err := fuzzyMock.VerifyUser(tokenMock)
	var expected error
	if err != expected {
		t.Fatalf("unexpected error, expected: %s, got: %s", expected, err)
	}

	if len(claims) != len(expectedClaims) {
		t.Fatalf("unexpected claims len, expected: %d, got: %d", len(expectedClaims), len(claims))
	}

	for k, v := range claims {
		if expectedClaims[k] != v {
			t.Fatalf("unexpected claim, expected: %s, got: %s", expected, err)
		}
	}
}
func Test_VerifyUser_Failed(t *testing.T) {
	tokenMock := "fake_token"
	expectedError := errors.New("unexpected error")
	jwtMock := jwtIssuerMock{
		validate: func(token string) (jwt.MapClaims, error) {
			if token == tokenMock {
				return nil, expectedError
			}
			return nil, nil
		},
	}

	fuzzyMock := core.NewFuzzy(nil, &jwtMock)

	_, err := fuzzyMock.VerifyUser(tokenMock)
	if !errors.Is(err, expectedError) {
		t.Fatalf("unexpected error, expected: %s, got: %s", expectedError, err)
	}
}
