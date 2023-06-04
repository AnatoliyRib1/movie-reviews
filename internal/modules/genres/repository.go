package genres

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

func (r *Repository) GetByID(ctx context.Context, id int) (*Genre, error) {
	var genre Genre
	query := "SELECT id, name FROM genres WHERE id = $1  "
	row := r.db.QueryRow(ctx, query, id)

	err := row.Scan(&genre.ID, &genre.Name)
	switch {
	case dbx.IsNoRows(err):
		return nil, errGenreWithIDNotFound(id)
	case err != nil:
		return nil, apperrors.Internal(err)

	}

	return &genre, nil
}

func (r *Repository) Create(ctx context.Context, genre *Genre) error {
	err := r.db.QueryRow(ctx, "insert into genres (name) values ($1) returning id", genre.Name).
		Scan(&genre.ID)
	switch {
	case dbx.IsUniqueViolation(err, "name"):
		return apperrors.AlreadyExists("genre", "name", genre.Name)
	case err != nil:
		return apperrors.Internal(err)

	}
	return nil
}

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

func errGenreWithIDNotFound(genreID int) error {
	return apperrors.NotFound("genre", "id", genreID)
}
