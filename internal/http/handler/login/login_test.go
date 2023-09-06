package login_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dgdraganov/fuzzy-user-api/internal/http/handler/login"
	"github.com/dgdraganov/fuzzy-user-api/internal/http/middleware"
	"github.com/dgdraganov/fuzzy-user-api/pkg/model"
	"go.uber.org/zap"
)

type loginMock struct {
	loginUser func(model.LoginDTO) (string, error)
}

func (r *loginMock) LoginUser(dto model.LoginDTO) (string, error) {
	return r.loginUser(dto)
}

func Test_ServeHTTP_Success(t *testing.T) {
	expectedToken := "fake_jwt_token"
	registry := &loginMock{
		loginUser: func(rd model.LoginDTO) (string, error) {
			return expectedToken, nil
		},
	}
	reg := login.NewLoginHandler(zap.NewNop().Sugar(), registry)
	loginHandler := middleware.SetContextRequestID(reg)
	dto := model.LoginDTO{
		Email:    "test@gmail.com",
		Password: "testPass",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(dto)
	if err != nil {
		t.Fatal("failed to encode RegisterDTO")
	}

	request, _ := http.NewRequest(http.MethodPost, "/api/register", &buf)
	response := httptest.NewRecorder()

	loginHandler.ServeHTTP(response, request)

	cookies := response.Result().Cookies()
	var authCookie string
	for _, cookie := range cookies {
		if cookie.Name == "Authentication" {
			authCookie = cookie.Value
			break
		}
	}

	respStruct := model.ResponseMessage{
		Message: "login successful",
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

	if authCookie != expectedToken {
		t.Fatalf("cookie does not match, expected: %s, got: %s", expectedToken, authCookie)
	}
	if string(expectedBody) != string(got) {
		t.Fatalf("response does not match, expected: %s, got: %s", expectedBody, got)
	}
	if expectedCode != response.Code {
		t.Fatalf("response code does not match, expected: %d, got: %d", expectedCode, response.Code)
	}
}

func Test_ServeHTTP_DeleteMethod_InvalidMethod(t *testing.T) {
	reg := login.NewLoginHandler(zap.NewNop().Sugar(), &loginMock{})
	loginHandler := middleware.SetContextRequestID(reg)

	testMethod := http.MethodDelete
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

func Test_ServeHTTP_InvalidBody(t *testing.T) {
	reg := login.NewLoginHandler(zap.NewNop().Sugar(), &loginMock{})
	loginHandler := middleware.SetContextRequestID(reg)

	body := strings.NewReader("{ invalid json }")

	request, _ := http.NewRequest(http.MethodPost, "/api/register", body)
	response := httptest.NewRecorder()

	loginHandler.ServeHTTP(response, request)

	respStruct := model.ResponseMessage{
		Message: "invalid request body",
	}
	expectedBody, err := json.Marshal(respStruct)
	if err != nil {
		t.Fatal("failed to marshal ResponseMessage")
	}
	expectedCode := http.StatusBadRequest
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
