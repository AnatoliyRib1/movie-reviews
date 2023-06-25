package movies

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/AnatoliyRib1/movie-reviews/internal/modules/stars"

	"github.com/AnatoliyRib1/movie-reviews/internal/modules/genres"

	"github.com/AnatoliyRib1/movie-reviews/internal/log"
)

type Service struct {
	repo         *Repository
	genreService *genres.Service
	starService  *stars.Service
}

func NewService(repo *Repository, genreService *genres.Service, starService *stars.Service) *Service {
	return &Service{
		repo:         repo,
		genreService: genreService,
		starService:  starService,
	}
}

func (s *Service) GetAllPaginated(ctx context.Context, searchTerm *string, starID *int, sortByRating *string, offset int, limit int) ([]*Movie, int, error) {
	return s.repo.GetAllPaginated(ctx, searchTerm, starID, sortByRating, offset, limit)
}

func (s *Service) Create(ctx context.Context, movie *MovieDetails) error {
	if err := s.repo.Create(ctx, movie); err != nil {
		return err
	}
	log.FromContext(ctx).Info("movie created", "movieID", movie.ID, "movieTitle", movie.Title)
	return s.assemble(ctx, movie)
}

func (s *Service) GetByID(ctx context.Context, movieID int) (*MovieDetails, error) {
	m, err := s.repo.GetByID(ctx, movieID)
	if err != nil {
		return nil, err
	}
	err = s.assemble(ctx, m)
	return m, err
}

func (s *Service) Update(ctx context.Context, movie *MovieDetails) error {
	if err := s.repo.Update(ctx, movie); err != nil {
		return err
	}
	log.FromContext(ctx).Info("movie updated", "movieTitle", movie.Title)
	return nil
}

func (s *Service) Delete(ctx context.Context, movieID int) error {
	if err := s.repo.Delete(ctx, movieID); err != nil {
		return err
	}
	log.FromContext(ctx).Info("movie deleted", "movieID", movieID)
	return nil
}

func (s *Service) assemble(ctx context.Context, movie *MovieDetails) error {
	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		var err error
		movie.Genres, err = s.genreService.GetByMovieID(groupCtx, movie.ID)
		return err
	})
	group.Go(func() error {
		var err error
		movie.Cast, err = s.starService.GetByMovieID(groupCtx, movie.ID)
		return err
	})

	return group.Wait()
}
