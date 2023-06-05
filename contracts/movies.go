package contracts

import "time"

type Movie struct {
	ID          int        `json:"id"`
	Title       string     `json:"title" `
	Description string     `json:"description"`
	ReleaseDate string     `json:"release_date"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
type GetMovieRequest struct {
	ID int `param:"movieId" validate:"nonzero"`
}

type GetMoviesRequest struct {
	PaginatedRequest
}

type CreateMovieRequest struct {
	Title       string    `json:"title" validate:"nonzero"`
	Description string    `json:"description" validate:"nonzero"`
	ReleaseDate time.Time `json:"release_date" validate:"nonzero"`
}

type UpdateMovieRequest struct {
	ID          int       `json:"id"`
	Title       string    `json:"title" validate:"nonzero"`
	Description string    `json:"description" validate:"nonzero"`
	ReleaseDate time.Time `json:"release_date" validate:"nonzero"`
}
type DeleteMovieRequest struct {
	ID int `param:"movieId" validate:"nonzero"`
}
