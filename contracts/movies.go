package contracts

import "time"

type Movie struct {
	ID          int        `json:"id"`
	Title       string     `json:"title" `
	ReleaseDate time.Time  `json:"release_date"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type MovieDetails struct {
	Movie
	Description string         `json:"description"`
	Version     int            `json:"version"`
	Genres      []*Genre       `json:"genres"`
	Cast        []*MovieCredit `json:"cast"`
}

type MovieCredit struct {
	Star    Star    `json:"star"`
	Role    string  `json:"role"`
	Details *string `json:"details,omitempty"`
}

type MovieCreditInfo struct {
	StarID  int     `json:"star_id"`
	Role    string  `json:"role"`
	Details *string `json:"details"`
}

type GetMovieRequest struct {
	MovieID int `param:"movieId" validate:"nonzero"`
}

type GetMoviesRequest struct {
	PaginatedRequest
	ID         *int
	SearchTerm *string
}

type CreateMovieRequest struct {
	Title       string             `json:"title" validate:"nonzero"`
	Description string             `json:"description" validate:"nonzero"`
	ReleaseDate time.Time          `json:"release_date" validate:"nonzero"`
	Genres      []int              `json:"genres"`
	Cast        []*MovieCreditInfo `json:"cast"`
}

type UpdateMovieRequest struct {
	MovieID     int                `json:"-" param:"movieId" validate:"nonzero"`
	Version     int                `json:"version" validate:"min=0"`
	Title       string             `json:"title" `
	Description string             `json:"description" `
	ReleaseDate time.Time          `json:"release_date" `
	Genres      []int              `json:"genres"`
	Cast        []*MovieCreditInfo `json:"cast"`
}
type DeleteMovieRequest struct {
	MovieID int `param:"movieId" validate:"nonzero"`
}
