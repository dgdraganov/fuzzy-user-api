package main

import (
	"fmt"
	"os"

	"github.com/dgdraganov/fuzzy-user-api/pkg/config"
)

func init() {
	config.LoadEnvConfig()
}

func main() {

	fmt.Println(os.Getenv("APP_ENV"))

}
