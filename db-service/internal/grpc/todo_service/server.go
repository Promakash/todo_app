package todo_service

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

func (s *DBService) Create(ctx context.Context, request *todov1.CreateRequest) (*todov1.CreateResponse, error) {
	if request.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "task's name is required")
	}
	if request.Description == "" {
		return nil, status.Error(codes.InvalidArgument, "task's description is required")
	}

	task := domain.Task{
		Name:        request.Name,
		Description: request.Description,
	}

	id, err := s.service.CreateTask(task)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &todov1.CreateResponse{Id: id}, nil
}

func (s *DBService) List(ctx context.Context, request *todov1.ListRequest) (*todov1.ListResponse, error) {
	tasks, err := s.service.ListTasks()
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	response := &todov1.ListResponse{}

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

func (s *DBService) Delete(ctx context.Context, request *todov1.DeleteRequest) (*todov1.DeleteResponse, error) {
	err := s.service.DeleteTask(request.Id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.InvalidArgument, "invalid task's id")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &todov1.DeleteResponse{}, nil
}

func (s *DBService) Done(ctx context.Context, request *todov1.DoneRequest) (*todov1.DoneResponse, error) {
	err := s.service.DoneTask(request.Id)

	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.InvalidArgument, "invalid task's id")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &todov1.DoneResponse{}, nil
}
