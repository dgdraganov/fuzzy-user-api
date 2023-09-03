package main

import (
	"github.com/dgdraganov/fuzzy-user-api/internal/http/server"
	"github.com/dgdraganov/fuzzy-user-api/pkg/config"
)

func init() {
	config.LoadEnvConfig()
}

func main() {
	fuzzy := server.NewHTTPServer()
	fuzzy.StartServer()
}
