package register

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgdraganov/fuzzy-user-api/pkg/model"
	"go.uber.org/zap"
)

type registerHandler struct {
	logs *zap.SugaredLogger
}

func NewRegisterHandler(logger *zap.SugaredLogger) *registerHandler {
	return &registerHandler{
		logs: logger,
	}
}

func (m *registerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		resp, err := json.Marshal(model.ErrorResponse{
			Error: "",
		})
		if err != nil {
			m.logs.Errorw(
				"failed to marshal error response",
			)
		}
		fmt.Println("delete me ")
		w.Write(resp)
	}
}
