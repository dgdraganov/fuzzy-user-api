package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadEnvConfig() {
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Errorf("loading configs: %w", err))
	}
}
