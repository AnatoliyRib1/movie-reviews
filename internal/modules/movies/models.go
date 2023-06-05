package movies

import "time"

type Movie struct {
	ID          int        `json:"id"`
	Title       string     `json:"title" `
	Description string     `json:"description"`
	ReleaseDate time.Time  `json:"release_date"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
