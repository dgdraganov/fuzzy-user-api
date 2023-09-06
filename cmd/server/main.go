package main

import (
	"github.com/dgdraganov/fuzzy-user-api/internal/http/server"
)

func main() {
	fuzzy := server.NewHTTPServer()
	fuzzy.StartServer()
}
