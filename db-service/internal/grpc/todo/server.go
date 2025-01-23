package todo

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"todo/db-service/internal/domain"
	"todo/db-service/internal/usecases"
	todov1 "todo/protos/gen/go"
)

type DBService struct {
	todov1.UnimplementedTodoServer
	service usecases.Task
}

func Register(gRPC *grpc.Server, service usecases.Task) {
	todov1.RegisterTodoServer(gRPC, &DBService{service: service})
}

func (s *DBService) CreateTask(ctx context.Context, request *todov1.CreateTaskRequest) (*todov1.CreateTaskResponse, error) {
	if request.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "task's name is required")
	}
	if request.GetDescription() == "" {
		return nil, status.Error(codes.InvalidArgument, "task's description is required")
	}

	task := domain.Task{
		Name:        request.GetName(),
		Description: request.GetDescription(),
	}

	id, err := s.service.CreateTask(task)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &todov1.CreateTaskResponse{Id: id}, nil
}

func (s *DBService) ListTasks(ctx context.Context, request *todov1.ListTasksRequest) (*todov1.ListTasksResponse, error) {
	tasks, err := s.service.ListTasks()
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	response := &todov1.ListTasksResponse{}

	for _, task := range tasks {
		response.Tasks = append(response.Tasks, &todov1.Task{
			Id:          task.ID,
			Name:        task.Name,
			Description: task.Description,
			IsDone:      task.IsDone,
		})
	}

	return response, nil
}

func (s *DBService) DeleteTaskByID(ctx context.Context, request *todov1.DeleteTaskByIDRequest) (*todov1.DeleteTaskByIDResponse, error) {
	err := s.service.DeleteTaskByID(request.Id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.InvalidArgument, "invalid task's id")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &todov1.DeleteTaskByIDResponse{}, nil
}

func (s *DBService) DoneTaskByID(ctx context.Context, request *todov1.DoneTaskByIDRequest) (*todov1.DoneTaskByIDResponse, error) {
	err := s.service.DoneTaskByID(request.Id)

	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.InvalidArgument, "invalid task's id")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &todov1.DoneTaskByIDResponse{}, nil
}

func (s *DBService) GetByID(ctx context.Context, request *todov1.GetByIDRequest) (*todov1.GetByIDResponse, error) {
	task, err := s.service.GetTaskByID(request.GetId())
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.InvalidArgument, "invalid task's id")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &todov1.GetByIDResponse{
		Task: &todov1.Task{
			Id:          task.ID,
			Name:        task.Name,
			Description: task.Description,
			IsDone:      task.IsDone,
		},
	}, nil
}
