package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dgdraganov/fuzzy-user-api/internal/http/handler/register"
	"github.com/dgdraganov/fuzzy-user-api/pkg/log"
	"github.com/dgdraganov/fuzzy-user-api/pkg/middleware"
	"github.com/dgdraganov/fuzzy-user-api/pkg/storage/pg"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type httpServer struct {
	mux      http.ServeMux
	register http.Handler
	logs     *zap.SugaredLogger
}

func NewHTTPServer() *httpServer {
	logger := log.NewZapLogger("fuzzy-user-api", zapcore.InfoLevel)

	conf := &pg.DbConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
		TimeZone: os.Getenv("DB_TIME_ZONE"),
	}

	db := pg.NewDatabase(conf)
	db.Connect()

	regHandler := register.NewRegisterHandler(logger, db)

	return &httpServer{
		mux:      *http.NewServeMux(),
		register: regHandler,
		logs:     logger,
	}
}

func (s *httpServer) RegisterHandlers() {
	// [POST]
	s.mux.Handle("/api/register", middleware.SetContextRequestID(s.register))
}

func (s *httpServer) StartServer() {
	port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	if err := http.ListenAndServe(port, nil); err != nil {
		s.logs.Fatalln("server stopped unexpectedly")
	}
}
