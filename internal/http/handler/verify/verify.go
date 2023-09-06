package verify

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dgdraganov/fuzzy-user-api/internal/http/common"
	"github.com/dgdraganov/fuzzy-user-api/pkg/jwt"
	"github.com/dgdraganov/fuzzy-user-api/pkg/model"
	"go.uber.org/zap"
)

type verifyHandler struct {
	logs     *zap.SugaredLogger
	registry Registry
}

func NewVerifyHandler(logger *zap.SugaredLogger, reg Registry) *verifyHandler {
	return &verifyHandler{
		logs:     logger,
		registry: reg,
	}
}

func (m *verifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(model.RequestID).(string)
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
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

	cookie, err := r.Cookie("Authentication")
	if err != nil {
		m.logs.Errorw(
			"failed to retrieve Auth cookie",
			"error", err,
			"request_id", requestID,
		)
		msg := "something went wrong on our end"
		status := http.StatusInternalServerError
		if errors.Is(err, http.ErrNoCookie) {
			msg = "could not find cookie 'Authentication'"
			status = http.StatusBadRequest
		}
		if err := common.WriteResponse(w, msg, status); err != nil {
			m.logs.Errorw(
				"write response failed (get cookie failed)",
				"error", err,
				"request_id", requestID,
			)
		}
		return
	}
	claims, err := m.registry.VerifyUser(cookie.Value)
	if err != nil {
		m.logs.Errorw(
			"failed verifying the user",
			"request_id", requestID,
			"error", err,
		)
		msg := "something went wrong on our end"
		status := http.StatusInternalServerError
		if errors.Is(err, jwt.ErrTokenNotValid) {
			msg = "invalid authentication token"
			status = http.StatusBadRequest
		}
		if err := common.WriteResponse(w, msg, status); err != nil {
			m.logs.Errorw(
				"write response failed (verifying user)",
				"error", err,
				"request_id", requestID,
			)
		}
		return
	}

	msg := fmt.Sprintf("user %s still logged in", claims["email"])
	if err := common.WriteResponse(w, msg, http.StatusOK); err != nil {
		m.logs.Errorw(
			"write response failed (verification success)",
			"error", err,
			"request_id", requestID,
		)
		return
	}

	m.logs.Infow(
		"successfully verified user",
		"request_id", requestID,
	)
}
