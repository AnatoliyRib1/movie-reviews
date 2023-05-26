package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v8"
)

type Config struct {
	DbUrl    string      `env:"DB_URL"`
	Port     int         `env:"PORT" envDefault:"8081"`
	Jwt      JwtConfig   `envPrefix:"JWT_"`
	Admin    AdminConfig `envPrefix:"ADMIN_"`
	Local    bool        `env:"LOCAL" envDefault:"false"`
	LogLevel string      `env:"LOG_LEVEL" envDefault:"info"`
}

type JwtConfig struct {
	Secret           string        `env:"SECRET"`
	AccessExpiration time.Duration `env:"ACCESS_EXPIRATION" envDefault:"30m"`
}

type AdminConfig struct {
	Username string `env:"NAME" validate:"min=5,max=16"`
	Email    string `env:"EMAIL" validate:"email"`
	Password string `env:"PASSWORD" validate:"password"`
}

func NewConfig() (*Config, error) {
	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &c, nil
}

func (ac *AdminConfig) IsSet() bool {
	return ac.Username != "" && ac.Email != "" && ac.Password != ""
}
