package common

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgdraganov/fuzzy-user-api/pkg/model"
)

func WriteResponse(w http.ResponseWriter, message string, statusCode int) error {
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
