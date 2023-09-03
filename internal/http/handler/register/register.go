package register

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dgdraganov/fuzzy-user-api/pkg/model"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type registerHandler struct {
	logs *zap.SugaredLogger
	repo UserRepository
}

func NewRegisterHandler(logger *zap.SugaredLogger, userRepo UserRepository) *registerHandler {
	return &registerHandler{
		logs: logger,
		repo: userRepo,
	}
}

func (m *registerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(model.RequestID).(string)
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		errMsg := model.ErrorResponse{
			Title: "user not registered",
			Error: fmt.Sprintf("invalid request method - %s", r.Method),
		}
		if err := writeResponse(w, errMsg, http.StatusMethodNotAllowed); err != nil {
			m.logs.Errorw(
				"error response failed",
				"error", err,
				"request_id", requestID,
			)
		}
		return
	}

	dto, err := parseRequestBody(r.Body)
	if err != nil {
		m.logs.Warnw(
			"parse request body failed",
			"error", err,
			"request_id", requestID,
		)
		errMsg := model.ErrorResponse{
			Title: "user not registered",
			Error: "bad request",
		}
		if err := writeResponse(w, errMsg, http.StatusBadRequest); err != nil {
			m.logs.Errorw(
				"error response failed",
				"error", err,
				"request_id", requestID,
			)
		}

		return
	}

	user, err := prepareUserStruct(dto)
	if err != nil {
		m.logs.Errorw(
			"prepare user struct failed",
			"error", err,
			"request_id", requestID,
		)
		errMsg := model.ErrorResponse{
			Title: "user not registered",
			Error: "something went kaput on our end",
		}
		if err := writeResponse(w, errMsg, http.StatusInternalServerError); err != nil {
			m.logs.Errorw(
				"error response failed",
				"error", err,
				"request_id", requestID,
			)
		}
		return
	}

	if err := m.repo.Create(&user); err != nil {
		m.logs.Errorw(
			"repo create user failed",
			"error", err,
			"request_id", requestID,
		)
		errMsg := model.ErrorResponse{
			Title: "user not registered",
			Error: "something went kaput on our end",
		}
		if err := writeResponse(w, errMsg, http.StatusInternalServerError); err != nil {
			m.logs.Errorw(
				"error response failed",
				"error", err,
				"request_id", requestID,
			)
		}
		return
	}

	m.logs.Infow(
		"user registered successfully",
		"request_id", requestID,
		"user_id", user.ID,
		"email", user.Email,
	)

	respMsg := model.SuccessResponse{
		Title: "user registered successfully",
	}
	if err := writeResponse(w, respMsg, http.StatusOK); err != nil {
		m.logs.Errorw(
			"error response failed",
			"error", err,
			"request_id", requestID,
		)
	}
}

func prepareUserStruct(dto model.RegisterDTO) (model.User, error) {
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

func parseRequestBody(b io.ReadCloser) (model.RegisterDTO, error) {

	var res model.RegisterDTO
	if err := json.NewDecoder(b).Decode(&res); err != nil {
		return model.RegisterDTO{}, fmt.Errorf("json decode: %w", err)
	}
	defer b.Close()
	validate := validator.New()
	err := validate.Struct(res)
	if err != nil {
		return model.RegisterDTO{}, fmt.Errorf("validate struct: %w", err)
	}
	return res, nil
}

func writeResponse(w http.ResponseWriter, errMsg any, statusCode int) error {
	resp, err := json.Marshal(errMsg)
	if err != nil {
		w.Write([]byte("Something went wrong!"))
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("json marshal: %w", err)
	}
	w.WriteHeader(statusCode)
	if _, err := w.Write(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("response write: %w", err)
	}
	return nil
}
