package repository

import (
	"context"
	"todo/db-service/domain"
)

type Task interface {
	PutTask(ctx context.Context, task domain.Task) (domain.TaskID, error)
	DeleteTaskByID(ctx context.Context, id domain.TaskID) error
	GetTaskByID(ctx context.Context, id domain.TaskID) (domain.Task, error)
	UpdateStatusByID(ctx context.Context, id domain.TaskID, status bool) error
	GetTasks(ctx context.Context) ([]domain.Task, error)
}
