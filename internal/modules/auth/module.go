package auth

import (
	"github.com/AnatoliyRib1/movie-reviews/internal/jwt"
	"github.com/AnatoliyRib1/movie-reviews/internal/modules/users"
)

type Module struct {
	Handler *Handler
	Service *Service
}

func NewModule(jwtService *jwt.Service, userService *users.Service) *Module {
	service := NewService(userService, jwtService)
	handler := NewHandler(service)

	return &Module{
		Handler: handler,
		Service: service,
	}
}
