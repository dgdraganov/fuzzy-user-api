package core_test

import (
	"errors"
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
				return model.User{}, gorm.ErrDuplicatedKey
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

/*func (f *fuzzy) LoginUser(dto model.LoginDTO) (string, error) {
err := validateLoginDTO(dto)
if err != nil {
	return "", fmt.Errorf("validate login dto: %w", err)
}

info, err := f.prepareTokenInfo(dto)
if err != nil {
	return "", fmt.Errorf("prapare token info: %w", err)
}

token := f.jwtIssuer.Generate(&info)
tokenStr, err := f.jwtIssuer.Sign(token)
if err != nil {
	return "", fmt.Errorf("token signing: %w", err)
}
return tokenStr, nil*/

func Test_LoginUser_Success(t *testing.T) {

}
