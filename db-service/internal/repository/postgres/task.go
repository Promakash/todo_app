package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"todo/db-service/internal/domain"
	"todo/db-service/internal/repository"
)

type TaskRepository struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) repository.Task {
	return &TaskRepository{
		pool: pool,
	}
}

func (r *TaskRepository) PutTask(ctx context.Context, task domain.Task) (domain.TaskID, error) {
	query := `INSERT INTO tasks (name, description, is_done)
			  VALUES ($1, $2, $3)
			  RETURNING id`

	var id domain.TaskID
	err := r.pool.QueryRow(ctx, query, task.Name, task.Description, task.IsDone).Scan(&id)
	return id, err
}

func (r *TaskRepository) DeleteTaskByID(ctx context.Context, id domain.TaskID) error {
	query := `DELETE FROM tasks
              WHERE id = $1`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *TaskRepository) UpdateStatusByID(ctx context.Context, id domain.TaskID, status bool) error {
	query := `UPDATE tasks
		      SET is_done  = $1
		      WHERE id = $2`

	cmdTag, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *TaskRepository) GetTasks(ctx context.Context) ([]domain.Task, error) {
	query := `SELECT id, name, description, is_done
			  FROM tasks`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.IsDone); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
