package verify_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgdraganov/fuzzy-user-api/internal/http/handler/verify"
	"github.com/dgdraganov/fuzzy-user-api/internal/http/middleware"
	"github.com/dgdraganov/fuzzy-user-api/pkg/model"
	"go.uber.org/zap"
)

type verifyMock struct {
	verifyUser func(jwtToken string) (map[string]any, error)
}

func (v *verifyMock) VerifyUser(jwtToken string) (map[string]any, error) {
	return v.verifyUser(jwtToken)
}

func Test_ServeHTTP_PostMethod_InvalidMethod(t *testing.T) {
	reg := verify.NewVerifyHandler(zap.NewNop().Sugar(), &verifyMock{})
	loginHandler := middleware.SetContextRequestID(reg)

	testMethod := http.MethodPost
	request, _ := http.NewRequest(testMethod, "/api/register", nil)
	response := httptest.NewRecorder()

	loginHandler.ServeHTTP(response, request)

	respStruct := model.ResponseMessage{
		Message: fmt.Sprintf("invalid request method - %s", testMethod),
	}
	expectedBody, err := json.Marshal(respStruct)
	if err != nil {
		t.Fatal("failed to marshal ResponseMessage")
	}
	got, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal("failed to read response body")
	}

	if string(expectedBody) != string(got) {
		t.Fatalf("response does not match, expected: %s, got: %s", expectedBody, got)
	}
	expectedCode := http.StatusMethodNotAllowed
	if expectedCode != response.Code {
		t.Fatalf("response code does not match, expected: %d, got: %d", expectedCode, response.Code)
	}
}

func Test_ServeHTTP_Success(t *testing.T) {
	testEmail := "test@test.com"
	registry := &verifyMock{
		verifyUser: func(jwtToken string) (map[string]any, error) {
			return map[string]any{
				"email": testEmail,
			}, nil
		},
	}
	reg := verify.NewVerifyHandler(zap.NewNop().Sugar(), registry)
	verifyHandler := middleware.SetContextRequestID(reg)

	request, _ := http.NewRequest(http.MethodGet, "/api/register", nil)
	response := httptest.NewRecorder()

	cookie := http.Cookie{}
	cookie.Name = "Authentication"
	cookie.Value = "fake_jwt_token"
	cookie.Expires = time.Now().Add(60 * time.Hour)
	cookie.Secure = false
	cookie.HttpOnly = true
	cookie.Path = "/"

	request.AddCookie(&cookie)

	verifyHandler.ServeHTTP(response, request)

	respStruct := model.ResponseMessage{
		Message: fmt.Sprintf("user %s still logged in", testEmail),
	}
	expectedBody, err := json.Marshal(respStruct)
	if err != nil {
		t.Fatal("failed to marshal ResponseMessage")
	}
	expectedCode := http.StatusOK
	got, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal("failed to read response body")
	}

	if string(expectedBody) != string(got) {
		t.Fatalf("response does not match, expected: %s, got: %s", expectedBody, got)
	}
	if expectedCode != response.Code {
		t.Fatalf("response code does not match, expected: %d, got: %d", expectedCode, response.Code)
	}
}
