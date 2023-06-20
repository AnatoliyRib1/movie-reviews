package genres

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

func (r *Repository) GetAll(ctx context.Context) ([]*Genre, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name FROM genres`)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	defer rows.Close()
	return pgx.CollectRows[*Genre](rows, pgx.RowToAddrOfStructByPos[Genre])
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

func (r *Repository) GetByMovieID(ctx context.Context, movieID int) ([]*Genre, error) {
	rows, err := r.db.Query(ctx, `
		select g.id, g.name from genres g
		inner join movie_genres mg on mg.genre_id = g.id
		where mg.movie_id = $1
		order by mg.order_no
		`, movieID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	defer rows.Close()
	return pgx.CollectRows[*Genre](rows, pgx.RowToAddrOfStructByPos[Genre])
}

func (r *Repository) GetRelationByMovieID(ctx context.Context, movieID int) ([]*MovieGenreRelation, error) {
	rows, err := dbx.FromContext(ctx, r.db).
		Query(ctx, "select movie_id, genre_id, order_no from movie_genres where movie_id = $1 order by order_no", movieID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	defer rows.Close()

	var relations []*MovieGenreRelation

	for rows.Next() {
		var relation MovieGenreRelation
		if err = rows.Scan(&relation.MovieID, &relation.GenreID, &relation.OrderNo); err != nil {
			return nil, apperrors.Internal(err)
		}
		relations = append(relations, &relation)
	}
	return relations, nil
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
