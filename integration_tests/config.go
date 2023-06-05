package tests

import (
	"time"

	"github.com/AnatoliyRib1/movie-reviews/internal/config"
)

const testPaginationSize = 2

func getConfig(pgConnString string) *config.Config {
	return &config.Config{
		DbURL: pgConnString,
		Port:  0,
		Jwt: config.JwtConfig{
			Secret:           "secret",
			AccessExpiration: time.Minute * 15,
		},
		Admin: config.AdminConfig{
			Username: "admin",
			Password: "&dm1Npa$$",
			Email:    "admin@email.com",
		},
		Pagination: config.PaginationConfig{
			DefaultSize: testPaginationSize,
			MaxSize:     50,
		},
		Local:    true,
		LogLevel: "error",
	}
}
