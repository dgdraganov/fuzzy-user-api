package register

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgdraganov/fuzzy-user-api/internal/http/common"
	"github.com/dgdraganov/fuzzy-user-api/pkg/model"
	"go.uber.org/zap"
)

type registerHandler struct {
	logs     *zap.SugaredLogger
	registry Registry
}

func NewRegisterHandler(logger *zap.SugaredLogger, reg Registry) *registerHandler {
	return &registerHandler{
		logs:     logger,
		registry: reg,
	}
}

func (m *registerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(model.RequestID).(string)
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		msg := fmt.Sprintf("invalid request method - %s", r.Method)
		if err := common.WriteResponse(w, msg, http.StatusMethodNotAllowed); err != nil {
			m.logs.Errorw(
				"write response failed (wrong request method)",
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
		if err := common.WriteResponse(w, "invalid request body", http.StatusBadRequest); err != nil {
			m.logs.Errorw(
				"write response failed (decode failed)",
				"error", err,
				"request_id", requestID,
			)
		}
		return
	}

	userExists, err := m.registry.UserExists(dto.Email)
	if err != nil {
		m.logs.Errorw(
			"failed getting user from db",
			"error", err,
			"request_id", requestID,
		)
		if err := common.WriteResponse(w, "something went wrong on our end", http.StatusInternalServerError); err != nil {
			m.logs.Errorw(
				"write response (get user failed)",
				"error", err,
				"request_id", requestID,
			)
		}
		return
	}

	if userExists {
		msg := fmt.Sprintf("user with email %s already exists", dto.Email)
		if err := common.WriteResponse(w, msg, http.StatusOK); err != nil {
			m.logs.Errorw(
				"write response (get user failed)",
				"error", err,
				"request_id", requestID,
			)
		}
		return
	}

	if err := m.registry.RegisterUser(dto); err != nil {
		m.logs.Errorw(
			"register user",
			"error", err,
			"request_id", requestID,
		)
		if err := common.WriteResponse(w, "invalid request body", http.StatusBadRequest); err != nil {
			m.logs.Errorw(
				"write response failed (failed register user)",
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

	if err := common.WriteResponse(w, "user registered successfully", http.StatusOK); err != nil {
		m.logs.Errorw(
			"write response (success)",
			"error", err,
			"request_id", requestID,
		)
	}
}
