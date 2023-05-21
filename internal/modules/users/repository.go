package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
		`insert into users (username, email, pass_hash) values ($1, $2, $3) returning id, created_at`,
		user.Username, user.Email, user.PasswordHash).
		Scan(&user.ID, &user.CreatedAt)
	return err
}

func (r *Repository) GetExistingUserWithPasswordByEmail(ctx context.Context, email string) (*UserWithPassword, error) {
	query := "SELECT id, username, email, pass_hash, role, bio FROM users WHERE email = $1 AND deleted_at IS NULL "
	row := r.db.QueryRow(ctx, query, email)

	user := newUserWithPassword()
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.Bio)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %s", err)
	}

	return user, nil
}

func (r *Repository) GetExistingUserById(ctx context.Context, userId int) (*User, error) {
	var user User
	query := "SELECT id, username, email,  role, bio FROM users WHERE email = $1 AND deleted_at IS NULL "
	row := r.db.QueryRow(ctx, query, userId)

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Bio)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %s", err)
	}

	return &user, nil
}

func (r *Repository) Delete(ctx context.Context, userId int) error {
	n, err := r.db.Exec(ctx, "UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL ", userId)
	if err != nil {
		return err
	}
	if n.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *Repository) Put(ctx context.Context, userId int, bio string) error {
	n, err := r.db.Exec(ctx, "UPDATE users SET bio = $2 WHERE id = $1 AND deleted_at IS NULL", userId, bio)
	if err != nil {
		return err
	}
	if n.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil

}
