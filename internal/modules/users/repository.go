package users

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
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
	query := "SELECT id, username, email, pass_hash, role, bio FROM users WHERE email = $1"
	row := r.db.QueryRow(ctx, query, email)

	user := UserWithPassword{
		User:         &User{},
		PasswordHash: "",
	}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.Bio)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %s", err)
	}

	return &user, nil
}

func (r *Repository) GetExistingUserById(ctx context.Context, userId int) (*User, error) {
	return nil, nil
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}

}
