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
	return s.repo.GetExistingUserWithPasswordByEmail(ctx, email)
}

func (s *Service) Delete(ctx context.Context, userId int) error {
	if err := s.repo.Delete(ctx, userId); err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateBio(ctx context.Context, userId int, bio string) error {
	if err := s.repo.UpdateBio(ctx, userId, bio); err != nil {
		return err
	}
	return nil
}

func (s *Service) Get(ctx context.Context, userId int) (user *User, err error) {
	return s.repo.GetExistingUserById(ctx, userId)
}

func (s *Service) SetRole(ctx context.Context, userId int, role string) error {
	if err := s.repo.SetRole(ctx, userId, role); err != nil {
		return err
	}
	return nil
}
