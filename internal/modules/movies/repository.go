package movies

import (
	"context"
	"time"

	"github.com/AnatoliyRib1/movie-reviews/internal/modules/stars"

	"github.com/AnatoliyRib1/movie-reviews/internal/modules/genres"
	"github.com/AnatoliyRib1/movie-reviews/internal/slices"

	"github.com/AnatoliyRib1/movie-reviews/internal/apperrors"
	"github.com/AnatoliyRib1/movie-reviews/internal/dbx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db        *pgxpool.Pool
	genreRepo *genres.Repository
	starRepo  *stars.Repository
}

func NewRepository(db *pgxpool.Pool, genreRepo *genres.Repository, starRepo *stars.Repository) *Repository {
	return &Repository{
		db:        db,
		genreRepo: genreRepo,
		starRepo:  starRepo,
	}
}

func (r *Repository) GetAllPaginated(ctx context.Context, offset int, limit int) ([]*Movie, int, error) {
	b := &pgx.Batch{}
	b.Queue("SELECT id, title,  release_date, created_at FROM movies WHERE deleted_at IS NULL ORDER BY id LIMIT $1 OFFSET $2", limit, offset)
	b.Queue("SELECT count(*) FROM movies WHERE deleted_at IS NULL")
	br := r.db.SendBatch(ctx, b)
	defer br.Close()

	rows, err := br.Query()
	if err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	defer rows.Close()

	var movies []*Movie
	for rows.Next() {
		var movie Movie
		if err = rows.Scan(&movie.ID, &movie.Title, &movie.ReleaseDate, &movie.CreatedAt); err != nil {
			return nil, 0, apperrors.Internal(err)
		}
		movies = append(movies, &movie)

	}
	if err = rows.Err(); err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	var total int
	if err = br.QueryRow().Scan(&total); err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	return movies, total, err
}

func (r *Repository) GetByID(ctx context.Context, id int) (*MovieDetails, error) {
	var movie MovieDetails
	query := "SELECT id, version ,title, description, release_date, created_at FROM movies WHERE id = $1 AND deleted_at IS NULL "
	row := r.db.QueryRow(ctx, query, id)

	err := row.Scan(&movie.ID, &movie.Version, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.CreatedAt)
	switch {
	case dbx.IsNoRows(err):
		return nil, errMovieWithIDNotFound(id)
	case err != nil:
		return nil, apperrors.Internal(err)

	}

	return &movie, nil
}

func (r *Repository) Create(ctx context.Context, movie *MovieDetails) error {
	err := dbx.InTransaction(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		err := tx.QueryRow(
			ctx,
			"insert into movies (title, description, release_date) values ($1, $2, $3) returning id, created_at",
			movie.Title, movie.Description, movie.ReleaseDate).
			Scan(&movie.ID, &movie.CreatedAt)
		if err != nil {
			return apperrors.Internal(err)
		}

		nextGenres := slices.MapIndex(movie.Genres, func(i int, g *genres.Genre) *genres.MovieGenreRelation {
			return &genres.MovieGenreRelation{
				MovieID: movie.ID,
				GenreID: g.ID,
				OrderNo: i,
			}
		})
		if err = r.updateGenres(ctx, nil, nextGenres); err != nil {
			return err
		}

		nextCast := slices.MapIndex(movie.Cast, func(i int, c *stars.MovieCredit) *stars.MovieStarRelation {
			return &stars.MovieStarRelation{
				MovieID: movie.ID,
				StarID:  c.Star.ID,
				Role:    c.Role,
				Details: c.Details,
				OrderNo: i,
			}
		})
		return r.updateCast(ctx, nil, nextCast)
	})
	if err != nil {
		return apperrors.Internal(err)
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, movie *MovieDetails) error {
	err := dbx.InTransaction(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		n, err := tx.Exec(ctx,
			"UPDATE movies SET version = version +1, title =$1, description = $2, release_date = $3 WHERE id = $4 AND version = $5 ",
			movie.Title, movie.Description, movie.ReleaseDate, movie.ID, movie.Version)
		if err != nil {
			return err
		}
		if n.RowsAffected() == 0 {
			_, err = r.GetByID(ctx, movie.ID)
			if err != nil {
				return err
			}
			return apperrors.VersionMismatch("movie", "id", movie.ID, movie.Version)
		}

		currentGenres, err := r.genreRepo.GetRelationByMovieID(ctx, movie.ID)
		if err != nil {
			return err
		}

		nextGenres := slices.MapIndex(movie.Genres, func(i int, g *genres.Genre) *genres.MovieGenreRelation {
			return &genres.MovieGenreRelation{
				GenreID: g.ID,
				MovieID: movie.ID,
				OrderNo: i,
			}
		})
		if err = r.updateGenres(ctx, currentGenres, nextGenres); err != nil {
			return err
		}

		currentCast, err := r.starRepo.GetRelationByMovieID(ctx, movie.ID)
		if err != nil {
			return err
		}

		nextCast := slices.MapIndex(movie.Cast, func(i int, c *stars.MovieCredit) *stars.MovieStarRelation {
			return &stars.MovieStarRelation{
				MovieID: movie.ID,
				StarID:  c.Star.ID,
				Role:    c.Role,
				Details: c.Details,
				OrderNo: i,
			}
		})
		if err = r.updateCast(ctx, currentCast, nextCast); err != nil {
			return err
		}
		return err
	})
	if err != nil {
		return apperrors.EnsureInternal(err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, movieID int) error {
	n, err := r.db.Exec(ctx, "UPDATE movies SET deleted_at = $2 WHERE id = $1 AND deleted_at IS NULL",
		movieID, time.Now())
	if err != nil {
		return apperrors.Internal(err)
	}
	if n.RowsAffected() == 0 {
		return errMovieWithIDNotFound(movieID)
	}
	return nil
}

func (r *Repository) updateGenres(ctx context.Context, current, next []*genres.MovieGenreRelation) error {
	q := dbx.FromContext(ctx, r.db)

	addFunc := func(mgo *genres.MovieGenreRelation) error {
		_, err := q.Exec(
			ctx,
			"insert into movie_genres (movie_id, genre_id, order_no) values ($1, $2, $3)",
			mgo.MovieID, mgo.GenreID, mgo.OrderNo)
		return err
	}

	removeFn := func(mgo *genres.MovieGenreRelation) error {
		_, err := q.Exec(
			ctx,
			"delete from movie_genres where movie_id = $1 and genre_id = $2",
			mgo.MovieID, mgo.GenreID)
		return err
	}
	return dbx.AdjustRelations(current, next, addFunc, removeFn)
}

func (r *Repository) updateCast(ctx context.Context, current, next []*stars.MovieStarRelation) error {
	q := dbx.FromContext(ctx, r.db)

	addFunc := func(mgo *stars.MovieStarRelation) error {
		_, err := q.Exec(
			ctx,
			"insert into movie_stars (movie_id, star_id,role, details, order_no) values ($1, $2, $3, $4, $5)",
			mgo.MovieID, mgo.StarID, mgo.Role, mgo.Details, mgo.OrderNo)
		return err
	}

	removeFn := func(mgo *stars.MovieStarRelation) error {
		_, err := q.Exec(
			ctx,
			"delete from movie_stars where movie_id = $1 and star_id = $2 and role = $3 and details = $4",
			mgo.MovieID, mgo.StarID, mgo.Role, mgo.Details)
		return err
	}
	return dbx.AdjustRelations(current, next, addFunc, removeFn)
}

func errMovieWithIDNotFound(movieID int) error {
	return apperrors.NotFound("movie", "id", movieID)
}
