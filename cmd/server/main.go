package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dgdraganov/fuzzy-user-api/internal/http/handler/register"
	"github.com/dgdraganov/fuzzy-user-api/pkg/config"
	"github.com/dgdraganov/fuzzy-user-api/pkg/log"
	"github.com/dgdraganov/fuzzy-user-api/pkg/middleware"
	"github.com/dgdraganov/fuzzy-user-api/pkg/storage/pg"
	"go.uber.org/zap/zapcore"
)

func init() {
	config.LoadEnvConfig()
}

func main() {

	//logger := log.NewZapLogger("fuzzy-user-api")

	logger := log.NewZapLogger("fuzzy-user-api", zapcore.InfoLevel)
	logger.Infow(
		"Service starting...",
		"app_env", os.Getenv("APP_ENV"),
		"app_port", os.Getenv("APP_PORT"),
	)

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

	rh := register.NewRegisterHandler(logger)

	// [POST]
	http.Handle("/api/register", middleware.SetContextRequestID(rh))

	// mux := http.NewServeMux()

	port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	if err := http.ListenAndServe(port, nil); err != nil {
		logger.Fatalln("")
	}
}
