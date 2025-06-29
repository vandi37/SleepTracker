package config

import (
	"github.com/goloop/env"
)

type Config struct {
	Port       int    `env:"PORT" def:"8080"`
	Refresh    JWT    `env:"REFRESH"`
	Access     JWT    `env:"ACCESS"`
	ConnString string `env:"CONN_STRING" def:"postgresql://user:password@localhost/app?sslmode=false"`
}

type JWT struct {
	Secret    string `env:"SECRET"`
	Expires   string `env:"EXP" def:"24h"`
	NotBefore string `env:"NBF" def:"1ms"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	if err := env.Unmarshal("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
