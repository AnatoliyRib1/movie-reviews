package genres

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

func (s *Service) GetAll(ctx context.Context) ([]*Genre, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) GetByID(ctx context.Context, genreID int) (genre *Genre, err error) {
	return s.repo.GetByID(ctx, genreID)
}

func (s *Service) Create(ctx context.Context, name string) (genre *Genre, err error) {
	genre = &Genre{Name: name}
	if err = s.repo.Create(ctx, genre); err != nil {
		return nil, err
	}
	log.FromContext(ctx).Info("genre created", "genreId", genre.ID, "genreName", genre.Name)
	return genre, nil
}

func (s *Service) Update(ctx context.Context, genreID int, name string) error {
	if err := s.repo.Update(ctx, genreID, name); err != nil {
		return err
	}
	log.FromContext(ctx).Info("genre updated", "genreId", genreID, "genreName", name)
	return nil
}

func (s *Service) Delete(ctx context.Context, genreID int) error {
	if err := s.repo.Delete(ctx, genreID); err != nil {
		return err
	}
	log.FromContext(ctx).Info("genre deleted", "genreId", genreID)
	return nil
}
