package users

import (
	"context"
)

type Service struct {
	repo *Repository
}

func (s *Service) Create(ctx context.Context, user *UserWithPassword) error {
	return s.repo.Create(ctx, user)

}

func (s *Service) GetUser(ctx context.Context, email string) (user *UserWithPassword, err error) {
	user, err = s.repo.GetExistingUserWithPasswordByEmail(ctx, email)
	return user, err

}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}

}
