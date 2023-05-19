package auth

import (
	"context"
	"fmt"
	"github.com/AnatoliyRib1/movie-reviews/internal/jwt"
	"github.com/AnatoliyRib1/movie-reviews/internal/modules/users"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userService *users.Service
	jwtService  jwt.Service
}

func (s *Service) Register(ctx context.Context, user *users.User, password string) error {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userWithPassword := &users.UserWithPassword{
		User:         user,
		PasswordHash: string(passHash),
	}
	return s.userService.Create(ctx, userWithPassword)

}

func (s *Service) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.userService.GetUser(ctx, email)
	if user == nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("incorrect password")
	}

	accessToken, err := s.jwtService.GenerateToken(int(user.ID), user.Role)
	return accessToken, nil
}

func NewService(userService *users.Service, jwtService *jwt.Service) *Service {
	return &Service{
		userService: userService,
		jwtService:  *jwtService,
	}

}
