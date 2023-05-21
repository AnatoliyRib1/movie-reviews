package users

import (
	"context"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, user *UserWithPassword) error {
	return s.repo.Create(ctx, user)
}

func (s *Service) GetUser(ctx context.Context, email string) (user *UserWithPassword, err error) {
	user, err = s.repo.GetExistingUserWithPasswordByEmail(ctx, email)
	return user, err
}

func (s *Service) Delete(ctx context.Context, userId int) error {
	return s.repo.Delete(ctx, userId)
}

func (s *Service) Update(ctx context.Context, userId int, bio string) error {
	return s.repo.Update(ctx, userId, bio)
}
