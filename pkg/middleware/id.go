package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dgdraganov/fuzzy-user-api/pkg/model"
	"github.com/google/uuid"
)

func SetContextRequestID(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New()
		ctx := context.WithValue(r.Context(), model.RequestID, requestID.String())
		r = r.WithContext(ctx)
		fmt.Println("==========================  Before request is handled  ==========================")

		handler.ServeHTTP(w, r)

		fmt.Println("==========================  After request is handled  ==========================")

	})
}
