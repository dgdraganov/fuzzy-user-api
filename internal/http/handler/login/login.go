package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgdraganov/fuzzy-user-api/internal/core"
	"github.com/dgdraganov/fuzzy-user-api/internal/http/common"
	"github.com/dgdraganov/fuzzy-user-api/pkg/model"
	"go.uber.org/zap"
)

type loginHandler struct {
	logs     *zap.SugaredLogger
	registry Registry
}

func NewRegisterHandler(logger *zap.SugaredLogger, reg Registry) *loginHandler {
	return &loginHandler{
		logs:     logger,
		registry: reg,
	}
}

func (m *loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	var dto model.LoginDTO
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

	token, err := m.registry.LoginUser(dto)
	if err != nil {
		m.logs.Errorw(
			"login user failed",
			"error", err,
			"request_id", requestID,
		)
		msg := "internal server error"
		status := http.StatusInternalServerError
		if errors.Is(err, core.ErrInvalidPassword) {
			msg = "incorrect password"
			status = http.StatusOK
		}
		if err := common.WriteResponse(w, msg, status); err != nil {
			m.logs.Errorw(
				"write response failed (login)",
				"error", err,
				"request_id", requestID,
			)
		}
		return
	}

	cookie := http.Cookie{}
	cookie.Name = "Authentication"
	cookie.Value = token
	cookie.Expires = time.Now().Add(60 * time.Hour)
	cookie.Secure = false
	cookie.HttpOnly = true
	cookie.Path = "/"

	http.SetCookie(w, &cookie)

	if err := common.WriteResponse(w, "login successful", http.StatusOK); err != nil {
		m.logs.Errorw(
			"write response failed (login success)",
			"error", err,
			"request_id", requestID,
		)
		return
	}

	m.logs.Infow(
		"successfully loged in",
		"jwt", token,
		"email", dto.Email,
		"request_id", requestID,
	)
}
