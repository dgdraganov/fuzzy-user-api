package login

import (
	"net/http"

	"go.uber.org/zap"
)

type loginHandler struct {
	logs *zap.SugaredLogger
	repo UserRepository
	jwt  JwtIssuer
}

func NewRegisterHandler(logger *zap.SugaredLogger, userRepo UserRepository, jwt JwtIssuer) *loginHandler {
	return &loginHandler{
		logs: logger,
		repo: userRepo,
		jwt:  jwt,
	}
}

func (m *loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")

}
