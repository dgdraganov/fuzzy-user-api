package register_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dgdraganov/fuzzy-user-api/internal/http/handler/register"
	"github.com/dgdraganov/fuzzy-user-api/internal/http/middleware"
	"github.com/dgdraganov/fuzzy-user-api/pkg/model"
	"go.uber.org/zap"
)

type registryMock struct {
	registerUser func(model.RegisterDTO) error
	userExists   func(email string) (bool, error)
}

func (r *registryMock) RegisterUser(dto model.RegisterDTO) error {
	return r.registerUser(dto)
}
func (r *registryMock) UserExists(email string) (bool, error) {
	return r.userExists(email)
}

func Test_ServeHTTP_Success(t *testing.T) {
	registry := &registryMock{
		registerUser: func(model.RegisterDTO) error {
			// simulate successful register
			return nil
		},
		userExists: func(email string) (bool, error) {
			return false, nil
		},
	}
	reg := register.NewRegisterHandler(zap.NewNop().Sugar(), registry)
	registerHandler := middleware.SetContextRequestID(reg)
	dto := model.RegisterDTO{
		FirstName: "Test",
		LastName:  "Test",
		Email:     "test@gmail.com",
		Password:  "testPass",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(dto)
	if err != nil {
		t.Fatal("failed to encode RegisterDTO")
	}

	request, _ := http.NewRequest(http.MethodPost, "/api/register", &buf)
	response := httptest.NewRecorder()

	registerHandler.ServeHTTP(response, request)

	respStruct := model.ResponseMessage{
		Message: "user registered successfully",
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

func Test_ServeHTTP_InvalidBody(t *testing.T) {
	registry := &registryMock{
		registerUser: func(model.RegisterDTO) error {
			// simulate successful register
			return nil
		},
		userExists: func(email string) (bool, error) {
			return false, nil
		},
	}
	reg := register.NewRegisterHandler(zap.NewNop().Sugar(), registry)
	registerHandler := middleware.SetContextRequestID(reg)

	body := strings.NewReader("{ invalid json }")

	request, _ := http.NewRequest(http.MethodPost, "/api/register", body)
	response := httptest.NewRecorder()

	registerHandler.ServeHTTP(response, request)

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

func Test_ServeHTTP_GetMethod_InvalidMethod(t *testing.T) {
	reg := register.NewRegisterHandler(zap.NewNop().Sugar(), &registryMock{})
	registerHandler := middleware.SetContextRequestID(reg)

	currMethod := http.MethodGet
	request, _ := http.NewRequest(currMethod, "/api/register", nil)
	response := httptest.NewRecorder()

	registerHandler.ServeHTTP(response, request)

	respStruct := model.ResponseMessage{
		Message: fmt.Sprintf("invalid request method - %s", currMethod),
	}
	expectedBody, err := json.Marshal(respStruct)
	if err != nil {
		t.Fatal("failed to marshal ResponseMessage")
	}
	expectedCode := http.StatusMethodNotAllowed
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

func Test_ServeHTTP_Duplicate(t *testing.T) {
	registry := &registryMock{
		registerUser: func(model.RegisterDTO) error {
			// simulate successful register
			return nil
		},
		userExists: func(email string) (bool, error) {
			return true, nil
		},
	}

	reg := register.NewRegisterHandler(zap.NewNop().Sugar(), registry)
	registerHandler := middleware.SetContextRequestID(reg)

	dto := model.RegisterDTO{
		FirstName: "Test",
		LastName:  "Test",
		Email:     "test@gmail.com",
		Password:  "testPass",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(dto)
	if err != nil {
		t.Fatal("failed to encode RegisterDTO")
	}

	request, _ := http.NewRequest(http.MethodPost, "/api/register", &buf)
	response := httptest.NewRecorder()

	registerHandler.ServeHTTP(response, request)

	respStruct := model.ResponseMessage{
		Message: fmt.Sprintf("user with email %s already exists", dto.Email),
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
