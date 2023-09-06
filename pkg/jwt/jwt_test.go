package jwt_test

import (
	"testing"
	"time"

	"github.com/dgdraganov/fuzzy-user-api/pkg/jwt"
	"github.com/dgdraganov/fuzzy-user-api/pkg/model"
)

func Test_Generate_Success(t *testing.T) {
	// Mock TimeNow from jwt package
	jwt.TimeNow = func() time.Time {
		time, err := time.Parse("2006-01-02", "2023-09-06")
		if err != nil {
			t.Fatal("failed to parse hardcoded date")
		}
		return time
	}

	jwtGen := jwt.NewJwtGenerator([]byte("test_secret"))
	data := model.TokenInfo{
		Email:      "test@gmail.com",
		FirstName:  "Test",
		LastName:   "Test",
		Subject:    "Login",
		Expiration: 2,
	}
	token := jwtGen.Generate(&data)
	tockernStr, err := jwtGen.Sign(token)
	if err != nil {
		t.Fatal("token singning failed")
	}

	expected := `eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAZ21haWwuY29tIiwiZXhwIjoxNjkzOTY1NjAwLCJmaXJzdF9uYW1lIjoiVGVzdCIsImlhdCI6MTY5Mzk1ODQwMCwibGFzdF9uYW1lIjoiVGVzdCIsInN1YiI6IkxvZ2luIn0.v9VXMcww_ngaIVXCvWq7yWv-DNfqXZ3q-rJP1nQJmYBcHk9cPGe6j2fot-xI_lJWqBbyWIRJoDjBeLaJgJY_Mw`
	if expected != tockernStr {
		t.Fatalf("token does not match, expected: %s, got: %s", expected, tockernStr)
	}
}
