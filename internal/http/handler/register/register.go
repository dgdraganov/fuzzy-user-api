package register

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgdraganov/fuzzy-user-api/pkg/model"
	"go.uber.org/zap"
)

type registerHandler struct {
	logs  *zap.SugaredLogger
	fuzzy Registry
	//repo UserRepository
}

func NewRegisterHandler(logger *zap.SugaredLogger) *registerHandler {
	return &registerHandler{
		logs: logger,
		//repo: userRepo,
	}
}

func (m *registerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(model.RequestID).(string)
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		msg := fmt.Sprintf("invalid request method - %s", r.Method)
		if err := writeResponse(w, msg, http.StatusMethodNotAllowed); err != nil {
			m.logs.Errorw(
				"error response failed",
				"error", err,
				"request_id", requestID,
			)
		}
		return
	}

	defer r.Body.Close()
	var dto model.RegisterDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		m.logs.Warnw(
			"json decode failed",
			"error", err,
			"request_id", requestID,
		)
		if err := writeResponse(w, "invalid request body", http.StatusBadRequest); err != nil {
			m.logs.Errorw(
				"write response failed",
				"error", err,
				"request_id", requestID,
			)
		}
		return
	}

	if err := m.fuzzy.RegisterUser(dto); err != nil {
		m.logs.Errorw(
			"register user",
			"error", err,
			"request_id", requestID,
		)
		if err := writeResponse(w, "invalid request body", http.StatusBadRequest); err != nil {
			m.logs.Errorw(
				"write response failed",
				"error", err,
				"request_id", requestID,
			)
		}
		return
	}

	m.logs.Infow(
		"user registered successfully",
		"request_id", requestID,
		"email", dto.Email,
	)

	if err := writeResponse(w, "user registered successfully", http.StatusOK); err != nil {
		m.logs.Errorw(
			"write response",
			"error", err,
			"request_id", requestID,
		)
	}
}

func writeResponse(w http.ResponseWriter, message string, statusCode int) error {
	respMsg := model.ResponseMessage{Message: message}
	resp, err := json.Marshal(respMsg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong!"))
		return fmt.Errorf("json marshal: %w", err)
	}
	w.WriteHeader(statusCode)
	if _, err := w.Write(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("response write: %w", err)
	}
	return nil
}
