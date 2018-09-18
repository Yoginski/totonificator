package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Token string
}

func Get() (*Config, error) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	errors := make([]string, 0, 1)
	if token == "" {
		errors = append(errors, "TELEGRAM_BOT_TOKEN environment variable not found")
	}
	if len(errors) > 0 {
		return nil, fmt.Errorf(strings.Join(errors, ", "))
	}
	return &Config{Token: token}, nil
}
