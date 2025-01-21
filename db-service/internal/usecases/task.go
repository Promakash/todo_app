package usecases

import "todo/db-service/internal/domain"

type Task interface {
	CreateTask(task domain.Task) (domain.TaskID, error)
	ListTasks() ([]domain.Task, error)
	DeleteTask(id domain.TaskID) error
	DoneTask(id domain.TaskID) error
}
