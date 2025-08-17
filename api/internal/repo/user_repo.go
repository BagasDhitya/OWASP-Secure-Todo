package repo

import (
	"context"
	"errors"

	"github.com/BagasDhitya/owasp-secure-todo/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct{ DB *pgxpool.Pool }

func (r *UserRepo) Create(ctx context.Context, u *models.User) error {
	q := `INSERT INTO users (username, email, password_hash, created_at, updated_at)
	      VALUES ($1, $2, $3, now(), now()) RETURNING id, created_at, updated_at`
	return r.DB.QueryRow(ctx, q, u.Username, u.Email, u.PasswordHash).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepo) ByEmail(ctx context.Context, email string) (*models.User, error) {
	q := `SELECT id, username, email, password_hash, created_at, updated_at
	      FROM users WHERE email=$1`
	row := r.DB.QueryRow(ctx, q, email)
	var u models.User
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
