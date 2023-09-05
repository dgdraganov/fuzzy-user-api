package core

import (
	"github.com/dgdraganov/fuzzy-user-api/pkg/model"
	"github.com/golang-jwt/jwt"
)

type Repository interface {
	Create(any) error
	GetUser(email string) (model.User, error)
}

type JwtIssuer interface {
	Generate(*model.TokenInfo) *jwt.Token
	Sign(token *jwt.Token) (string, error)
}
