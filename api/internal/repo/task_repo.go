package repo

import (
	"context"
	"errors"

	"github.com/BagasDhitya/owasp-secure-todo/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("record not found")

type TaskRepo struct{ DB *pgxpool.Pool }

func (r *TaskRepo) ListByUser(ctx context.Context, userID int64) ([]models.Task, error) {
	q := `SELECT id, title, description, status, created_at, updated_at
	      FROM tasks WHERE user_id=$1 ORDER BY created_at DESC`
	rows, err := r.DB.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (r *TaskRepo) Create(ctx context.Context, userID int64, t *models.Task) error {
	q := `INSERT INTO tasks (user_id, title, description, status)
	      VALUES ($1,$2,$3,COALESCE($4,'pending'))
	      RETURNING id, created_at, updated_at`
	return r.DB.QueryRow(ctx, q, userID, t.Title, t.Description, t.Status).
		Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *TaskRepo) Update(ctx context.Context, userID, id int64, t *models.Task) error {
	q := `UPDATE tasks SET title=$1, description=$2, status=$3, updated_at=now()
	      WHERE id=$4 AND user_id=$5 RETURNING updated_at`
	err := r.DB.QueryRow(ctx, q, t.Title, t.Description, t.Status, id, userID).Scan(&t.UpdatedAt)
	if err != nil {
		return ErrNotFound
	}
	return nil
}

func (r *TaskRepo) Delete(ctx context.Context, userID, id int64) error {
	cmd, err := r.DB.Exec(ctx, `DELETE FROM tasks WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
