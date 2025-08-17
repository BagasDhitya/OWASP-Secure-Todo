package repo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgRefreshStore struct {
	DB *pgxpool.Pool
}

func (s *PgRefreshStore) Save(ctx context.Context, userID int64, jti string, rawToken string, exp time.Time, ua, ip string) error {
	_, err := s.DB.Exec(ctx,
		`INSERT INTO refresh_tokens (user_id, jti, token, expires_at, user_agent, ip_address)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		userID, jti, rawToken, exp, ua, ip,
	)
	return err
}

func (s *PgRefreshStore) Revoke(ctx context.Context, jti string) error {
	_, err := s.DB.Exec(ctx, `DELETE FROM refresh_tokens WHERE jti=$1`, jti)
	return err
}

func (s *PgRefreshStore) FindValid(ctx context.Context, userID int64, jti, rawToken string) (bool, error) {
	var count int
	err := s.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM refresh_tokens
		 WHERE user_id=$1 AND jti=$2 AND token=$3 AND expires_at > NOW()`,
		userID, jti, rawToken,
	).Scan(&count)
	return count > 0, err
}
