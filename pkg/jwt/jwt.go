package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgdraganov/fuzzy-user-api/pkg/model"
	"github.com/golang-jwt/jwt"
)

var TimeNow = time.Now
var ErrTokenNotValid error = errors.New("token is not valid")

type jwtGenerator struct {
	secret []byte
}

func NewJwtGenerator(jwtSecret []byte) *jwtGenerator {
	return &jwtGenerator{
		secret: jwtSecret,
	}
}

func (gen *jwtGenerator) Generate(data *model.TokenInfo) *jwt.Token {
	token := jwt.New(jwt.SigningMethodHS512)

	token.Header["typ"] = "JWT"
	token.Header["alg"] = jwt.SigningMethodHS512.Name

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = data.Subject
	claims["iat"] = TimeNow().Unix()
	claims["exp"] = TimeNow().Unix() + int64(data.Expiration*3600)
	claims["first_name"] = data.FirstName
	claims["last_name"] = data.LastName
	claims["email"] = data.Email

	return token
}

func (gen *jwtGenerator) Sign(token *jwt.Token) (string, error) {
	tokenStr, err := token.SignedString(gen.secret)
	if err != nil {
		return "", fmt.Errorf("string signed: %w", err)
	}

	return tokenStr, nil
}

func (gen *jwtGenerator) Validate(token string) (jwt.MapClaims, error) {

	fmt.Println(token)
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return gen.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwt parse: %w", err)
	}

	if !jwtToken.Valid {
		return nil, ErrTokenNotValid
	}
	var claims jwt.MapClaims
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("jwt claims type assertion failed")
	}
	return claims, nil
}
