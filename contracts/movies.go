package contracts

import (
	"strconv"
	"time"
)

type Movie struct {
	ID          int        `json:"id"`
	Title       string     `json:"title" `
	ReleaseDate time.Time  `json:"release_date"`
	AvgRating   *float64   `json:"avg_rating,omitempty"`
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
	StarID       *int    `query:"starId"`
	SearchTerm   *string `query:"q"`
	SortByRating *string `json:"sortByRating" validate:"sort"`
}

func (r *GetMoviesRequest) ToQueryParams() map[string]string {
	param := r.PaginatedRequest.ToQueryParams()
	if r.StarID != nil {
		param["starId"] = strconv.Itoa(*r.StarID)
	}
	if r.SearchTerm != nil {
		param["q"] = *r.SearchTerm
	}
	if r.SortByRating != nil {
		param["sortByRating"] = *r.SortByRating
	}
	return param
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
