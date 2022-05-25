package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	Address          string `env:"ADDRESS"`
	ConnectionString string `env:"CONNECTION_STRING,required"`
	SecretKey        string `env:"SECRET_KEY,unset"`
	StaticBase       string `env:"STATIC_BASE"`
	PasswordCost     int    `env:"PASSWORD_COST"`
	Production       bool   `env:"PRODUCTION"`
	// DBTimeout        int    `env:"DB_TIMEOUT"`
}

func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	cfg := Config{}
	err = env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
