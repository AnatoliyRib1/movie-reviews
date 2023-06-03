package users

import (
	"context"

	"github.com/AnatoliyRib1/movie-reviews/internal/log"
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

func (s *Service) GetExistingUserWithPasswordByEmail(ctx context.Context, email string) (user *UserWithPassword, err error) {
	return s.repo.GetExistingUserWithPasswordByEmail(ctx, email)
}

func (s *Service) Delete(ctx context.Context, userID int) error {
	if err := s.repo.Delete(ctx, userID); err != nil {
		return err
	}
	log.FromContext(ctx).Info("user deleted", "userId")
	return nil
}

func (s *Service) Update(ctx context.Context, userID int, bio string) error {
	if err := s.repo.Update(ctx, userID, bio); err != nil {
		return err
	}
	log.FromContext(ctx).Info("user updated", "userId")
	return nil
}

func (s *Service) GetExistingUserByID(ctx context.Context, userID int) (user *User, err error) {
	return s.repo.GetExistingUserByID(ctx, userID)
}

func (s *Service) SetRole(ctx context.Context, userID int, role string) error {
	if err := s.repo.SetRole(ctx, userID, role); err != nil {
		return err
	}
	log.FromContext(ctx).Info("user role updated", "userId", userID, "role", role)
	return nil
}

func (s *Service) GetExistingUserByUserName(ctx context.Context, userName string) (user *User, err error) {
	return s.repo.GetExistingUserByUserName(ctx, userName)
}
