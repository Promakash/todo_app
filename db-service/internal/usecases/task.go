package usecases

import "todo/db-service/internal/domain"

type Task interface {
	CreateTask(task domain.Task) (domain.TaskID, error)
	ListTasks() ([]domain.Task, error)
	DeleteTaskByID(id domain.TaskID) error
	DoneTaskByID(id domain.TaskID) error
	GetTaskByID(id domain.TaskID) (domain.Task, error)
}
