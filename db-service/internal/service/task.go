package service

import (
	"context"
	"log/slog"
	"todo/db-service/internal/domain"
	"todo/db-service/internal/repository"
	"todo/db-service/internal/usecases"
	pkglog "todo/pkg/log"
)

type TaskService struct {
	log  *slog.Logger
	repo repository.Task
}

func NewTaskService(repo repository.Task, log *slog.Logger) usecases.Task {
	return &TaskService{
		repo: repo,
		log:  log,
	}
}

func (s *TaskService) CreateTask(task domain.Task) (domain.TaskID, error) {
	const op = "Task.CreateTask"

	log := s.log.With(
		slog.String("op", op),
		slog.String("Name", task.Name),
	)

	log.Info("Trying to add new task")

	id, err := s.repo.PutTask(context.Background(), task)
	if err != nil {
		log.Error("can't create new task: ", pkglog.Err(err))
	}

	return id, err
}
func (s *TaskService) ListTasks() ([]domain.Task, error) {
	const op = "Task.ListTasks"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("Trying to list tasks")

	tasks, err := s.repo.GetTasks(context.Background())
	if err != nil {
		log.Error("can't list tasks: ", pkglog.Err(err))
	}

	return tasks, err
}
func (s *TaskService) DeleteTask(id domain.TaskID) error {
	const op = "Task.ListTasks"

	log := s.log.With(
		slog.String("op", op),
		slog.Uint64("id", id),
	)

	log.Info("Trying to delete task")

	err := s.repo.DeleteTaskByID(context.Background(), id)
	if err != nil {
		log.Error("can't delete task: ", pkglog.Err(err))
	}

	return err
}
func (s *TaskService) DoneTask(id domain.TaskID) error {
	const op = "Task.DoneTask"

	log := s.log.With(
		slog.String("op", op),
		slog.Uint64("id", id),
	)

	log.Info("Trying to mark task as done")

	err := s.repo.UpdateStatusByID(context.Background(), id, true)
	if err != nil {
		log.Error("can't mark task as done: ", pkglog.Err(err))
	}

	return err
}
