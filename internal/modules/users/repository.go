package users

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

func (r *Repository) Create(ctx context.Context, user *UserWithPassword) error {
	err := r.db.QueryRow(
		ctx,
		`insert into users (username, email, pass_hash, role) values ($1, $2, $3, $4) returning id, created_at`,
		user.Username, user.Email, user.PasswordHash, user.Role).
		Scan(&user.ID, &user.CreatedAt)
	switch {
	case dbx.IsUniqueViolation(err, "email"):
		return apperrors.AlreadyExists("user", "email", user.Email)
	case dbx.IsUniqueViolation(err, "username"):
		return apperrors.AlreadyExists("user", "email", user.Username)
	case err != nil:
		return apperrors.Internal(err)

	}
	return nil
}

func (r *Repository) GetExistingUserWithPasswordByEmail(ctx context.Context, email string) (*UserWithPassword, error) {
	query := "SELECT id, username, email, pass_hash, role, bio FROM users WHERE email = $1 AND deleted_at IS NULL "
	row := r.db.QueryRow(ctx, query, email)

	user := newUserWithPassword()
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.Bio)
	switch {
	case dbx.IsNoRows(err):
		return nil, errUserWithEmailNotFound(user.Email)
	case err != nil:
		return nil, apperrors.Internal(err)

	}
	return user, nil
}

func (r *Repository) GetExistingUserById(ctx context.Context, userId int) (*User, error) {
	var user User
	query := "SELECT id, username, email,  role, bio FROM users WHERE id = $1 AND deleted_at IS NULL "
	row := r.db.QueryRow(ctx, query, userId)

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Bio)
	switch {
	case dbx.IsNoRows(err):
		return nil, errUserWithIdNotFound(userId)
	case err != nil:
		return nil, apperrors.Internal(err)

	}

	return &user, nil
}

func (r *Repository) GetExistingUserByUserName(ctx context.Context, userName string) (*User, error) {
	var user User
	query := "SELECT id, username, email,  role, bio FROM users WHERE username = $1 AND deleted_at IS NULL "
	row := r.db.QueryRow(ctx, query, userName)

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Bio)
	switch {
	case dbx.IsNoRows(err):
		return nil, errUserWithUserNameNotFound(userName)
	case err != nil:
		return nil, apperrors.Internal(err)

	}

	return &user, nil
}

func (r *Repository) Delete(ctx context.Context, userId int) error {
	n, err := r.db.Exec(ctx, "UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL ", userId)
	if err != nil {
		return apperrors.Internal(err)
	}
	if n.RowsAffected() == 0 {
		return errUserWithIdNotFound(userId)
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, userId int, bio string) error {
	n, err := r.db.Exec(ctx, "UPDATE users SET bio = $2 WHERE id = $1 AND deleted_at IS NULL", userId, bio)
	if err != nil {
		return apperrors.Internal(err)
	}
	if n.RowsAffected() == 0 {
		return errUserWithIdNotFound(userId)
	}
	return nil
}

func (r *Repository) SetRole(ctx context.Context, userId int, role string) error {
	n, err := r.db.Exec(ctx, "UPDATE users SET role = $2 WHERE id = $1 AND deleted_at IS NULL", userId, role)
	if err != nil {
		return apperrors.Internal(err)
	}
	if n.RowsAffected() == 0 {
		return errUserWithIdNotFound(userId)
	}
	return nil
}

func errUserWithIdNotFound(userId int) error {
	return apperrors.NotFound("user", "id", userId)
}

func errUserWithUserNameNotFound(userName string) error {
	return apperrors.NotFound("user", "username", userName)
}

func errUserWithEmailNotFound(userEmail string) error {
	return apperrors.NotFound("user", "email", userEmail)
}
