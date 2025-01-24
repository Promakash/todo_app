package grpc

import (
	"context"
	"fmt"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"
	apidomain "todo/api-gateway/internal/domain"
	"todo/db-service/domain"
	pkggrpc "todo/pkg/grpc"
	todov1 "todo/protos/gen/go"
)

type Client struct {
	api todov1.TodoClient
	log *slog.Logger
}

func New(ctx context.Context,
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	const op = "grpc.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(pkggrpc.InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.With(slog.String("op", op)).Info("client has started")

	return &Client{
		api: todov1.NewTodoClient(cc),
		log: log,
	}, nil
}

func (c *Client) ListTasks(ctx context.Context) ([]domain.Task, error) {
	req := &todov1.ListTasksRequest{}
	resp, err := c.api.ListTasks(ctx, req)
	if err != nil {
		return nil, apidomain.HandleGRPCError(err)
	}

	tasks := make([]domain.Task, len(resp.GetTasks()))
	for i, task := range resp.GetTasks() {
		tasks[i] = domain.TaskFromGRPC(task)
	}

	return tasks, nil
}

func (c *Client) GetTaskByID(ctx context.Context, id domain.TaskID) (domain.Task, error) {
	req := &todov1.GetByIDRequest{Id: id}
	var task domain.Task

	resp, err := c.api.GetByID(ctx, req)
	if err != nil {
		return task, apidomain.HandleGRPCError(err)
	}
	task = domain.TaskFromGRPC(resp.GetTask())

	return task, nil
}

func (c *Client) CreateTask(ctx context.Context, name, description string) (domain.TaskID, error) {
	req := &todov1.CreateTaskRequest{
		Name:        name,
		Description: description,
	}
	var id domain.TaskID

	resp, err := c.api.CreateTask(ctx, req)
	if err != nil {
		return id, apidomain.HandleGRPCError(err)
	}
	id = resp.GetId()

	return id, nil
}

func (c *Client) DeleteTaskByID(ctx context.Context, id domain.TaskID) error {
	req := &todov1.DeleteTaskByIDRequest{Id: id}

	_, err := c.api.DeleteTaskByID(ctx, req)

	return apidomain.HandleGRPCError(err)
}

func (c *Client) DoneTaskByID(ctx context.Context, id domain.TaskID) error {
	req := &todov1.DoneTaskByIDRequest{Id: id}

	_, err := c.api.DoneTaskByID(ctx, req)

	return apidomain.HandleGRPCError(err)
}
