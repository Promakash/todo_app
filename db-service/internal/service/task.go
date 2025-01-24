package service

import (
	"context"
	"log/slog"
	"time"
	"todo/db-service/domain"
	"todo/db-service/internal/repository"
	"todo/db-service/internal/usecases"
	"todo/pkg/infra/cache"
	pkglog "todo/pkg/log"
)

type TaskService struct {
	log        *slog.Logger
	repo       repository.Task
	cache      cache.Cache
	defaultTTL time.Duration
}

func NewTaskService(log *slog.Logger, repo repository.Task, cache cache.Cache, defaultTTL time.Duration) usecases.Task {
	return &TaskService{
		repo:       repo,
		log:        log,
		cache:      cache,
		defaultTTL: defaultTTL,
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
		return nil, err
	}

	return tasks, err
}

func (s *TaskService) DeleteTaskByID(id domain.TaskID) error {
	const op = "Task.DeleteTaskByID"
	ctx := context.Background()

	log := s.log.With(
		slog.String("op", op),
		slog.Uint64("id", id),
	)

	log.Info("Trying to delete task")

	err := s.repo.DeleteTaskByID(ctx, id)
	if err != nil {
		log.Error("can't delete task: ", pkglog.Err(err))
	} else {
		go func() { _ = s.cache.Delete(ctx, domain.TaskIDToString(id)) }()
	}

	return err
}

func (s *TaskService) DoneTaskByID(id domain.TaskID) error {
	const op = "Task.DoneTaskByID"
	ctx := context.Background()

	log := s.log.With(
		slog.String("op", op),
		slog.Uint64("id", id),
	)

	log.Info("Trying to mark task as done")

	err := s.repo.UpdateStatusByID(ctx, id, true)
	if err != nil {
		log.Error("can't mark task as done: ", pkglog.Err(err))
	} else {
		go func() { _ = s.cache.Delete(ctx, domain.TaskIDToString(id)) }()
	}

	return err
}

func (s *TaskService) GetTaskByID(id domain.TaskID) (domain.Task, error) {
	const op = "Task.GetTaskByID"
	ctx := context.Background()

	log := s.log.With(
		slog.String("op", op),
		slog.Uint64("id", id),
	)

	log.Info("Trying to get task by id")

	var task domain.Task

	err := s.cache.Get(ctx, domain.TaskIDToString(id), &task)
	if err == nil {
		return task, nil
	} else {
		log.Error("error while getting value from cache: ", pkglog.Err(err))
	}

	task, err = s.repo.GetTaskByID(ctx, id)
	if err != nil {
		log.Error("can't get task by id: ", pkglog.Err(err))
	} else {
		_ = s.cache.Set(ctx, domain.TaskIDToString(id), task, s.defaultTTL)
	}

	return task, err
}
