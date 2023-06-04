package stars

import (
	"context"

	"github.com/jackc/pgx/v5"

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

func (r *Repository) GetAllPaginated(ctx context.Context, offset int, limit int) ([]*Star, int, error) {
	b := &pgx.Batch{}
	b.Queue("SELECT id, first_name ,middle_name, last_name, birth_date, birth_place, death_date, bio, created_at, deleted_at FROM stars WHERE deleted_at IS NULL ORDER BY id LIMIT $1 OFFSET $2", limit, offset)
	b.Queue("SELECT count(*) FROM stars WHERE deleted_at IS NULL")
	br := r.db.SendBatch(ctx, b)
	defer br.Close()

	rows, err := br.Query()
	if err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	defer rows.Close()

	var stars []*Star
	for rows.Next() {
		var star Star
		if err = rows.Scan(&star.ID, &star.FirstName, &star.MiddleName, &star.LastName, &star.BirthDate, &star.BirthPlace, &star.DeathDate, &star.Bio, &star.CreatedAt, &star.DeletedAt); err != nil {
			return nil, 0, apperrors.Internal(err)
		}
		stars = append(stars, &star)

	}
	if err = rows.Err(); err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	var total int
	if err = br.QueryRow().Scan(&total); err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	return stars, total, err
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

	err := row.Scan(&star.ID, &star.FirstName, &star.MiddleName, &star.LastName, &star.BirthDate, &star.BirthPlace, &star.DeathDate, &star.Bio, &star.CreatedAt)
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
		"insert into stars (first_name ,middle_name, last_name, birth_date, birth_place, death_date, bio) values ($1, $2, $3, $4, $5, $6, $7) returning id, created_at, deleted_at",
		star.FirstName, star.MiddleName, star.LastName, star.BirthDate, star.BirthPlace, star.DeathDate, star.Bio).
		Scan(&star.ID, &star.CreatedAt, &star.DeletedAt)
	if err != nil {
		return apperrors.Internal(err)
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, star *Star) error {
	n, err := r.db.Exec(ctx, "UPDATE stars SET first_name =$1, middle_name = $2, last_name = $3, birth_date = $4, birth_place = $5, death_date = $6, bio = $7 WHERE id = $8 ", star.FirstName, star.MiddleName, star.LastName, star.BirthDate, star.BirthPlace, star.DeathDate, star.Bio, star.ID)
	if err != nil {
		return apperrors.Internal(err)
	}
	if n.RowsAffected() == 0 {
		return errStarWithIDNotFound(star.ID)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, starID int) error {
	n, err := r.db.Exec(ctx, "DELETE FROM stars WHERE id = $1 ", starID)
	if err != nil {
		return apperrors.Internal(err)
	}
	if n.RowsAffected() == 0 {
		return errStarWithIDNotFound(starID)
	}
	return nil
}

func errStarWithIDNotFound(starID int) error {
	return apperrors.NotFound("star", "id", starID)
}
