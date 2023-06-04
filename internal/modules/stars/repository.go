package stars

import (
	"context"

	"github.com/AnatoliyRib1/movie-reviews/internal/apperrors"
	"github.com/AnatoliyRib1/movie-reviews/internal/dbx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

/*
	func (r *Repository) GetAll(ctx context.Context) ([]*Genre, error) {
		rows, err := r.db.Query(ctx, `SELECT id, name FROM genres`)
		if err != nil {
			return nil, apperrors.Internal(err)
		}
		var genres []*Genre
		defer rows.Close()

		for rows.Next() {
			var genre Genre
			if err = rows.Scan(&genre.ID, &genre.Name); err != nil {
				return nil, apperrors.Internal(err)
			}
			genres = append(genres, &genre)

		}
		if err = rows.Err(); err != nil {
			return nil, apperrors.Internal(err)
		}
		return genres, err
	}
*/
func (r *Repository) GetByID(ctx context.Context, id int) (*Star, error) {
	var star Star
	query := "SELECT id, first_name ,middle_name, last_name, birth_date,birth_place, death_date, bio, created_at FROM stars WHERE id = $1"
	row := r.db.QueryRow(ctx, query, id)

	err := row.Scan(&star.ID, &star.FirstName, &star.MiddleName, &star.LastName, &star.BirthDate, &star.BirthPlace, &star.DeathDate, &star.Bio, &star.CreatedAd)
	switch {
	case dbx.IsNoRows(err):
		return nil, errStarWithIDNotFound(id)
	case err != nil:
		return nil, apperrors.Internal(err)

	}

	return &star, nil
}

func (r *Repository) Create(ctx context.Context, star *Star) error {
	err := r.db.QueryRow(ctx,
		"insert into stars (first_name ,middle_name, last_name, birth_date, birth_place, death_date, bio) values ($1, $2, $3, $4, $5, $6, $7) returning id, created_at",
		star.FirstName, star.MiddleName, star.LastName, star.BirthDate, star.BirthPlace, star.DeathDate, star.Bio).
		Scan(&star.ID, &star.CreatedAd)
	if err != nil {
		return apperrors.Internal(err)
	}
	return nil
}

/*
	func (r *Repository) Update(ctx context.Context, genreID int, name string) error {
		n, err := r.db.Exec(ctx, "UPDATE genres SET name = $2 WHERE id = $1 ", genreID, name)
		if err != nil {
			return apperrors.Internal(err)
		}
		if n.RowsAffected() == 0 {
			return errGenreWithIDNotFound(genreID)
		}
		return nil
	}

	func (r *Repository) Delete(ctx context.Context, genreID int) error {
		n, err := r.db.Exec(ctx, "DELETE FROM genres WHERE id = $1 ", genreID)
		if err != nil {
			return apperrors.Internal(err)
		}
		if n.RowsAffected() == 0 {
			return errGenreWithIDNotFound(genreID)
		}
		return nil
	}
*/
func errStarWithIDNotFound(starID int) error {
	return apperrors.NotFound("star", "id", starID)
}
