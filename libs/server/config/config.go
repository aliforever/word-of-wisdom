package config

import (
	"errors"
	"os"
)

var (
	emptyAddressErr = errors.New("empty_address")
)

type Config struct {
	Address string
}

func LoadFromEnv() (*Config, error) {
	address := os.Getenv("ADDRESS")
	if address == "" {
		return nil, emptyAddressErr
	}

	return &Config{Address: address}, nil
}
