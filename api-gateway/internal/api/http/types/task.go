package types

import (
	"net/http"
	apidomain "todo/api-gateway/internal/domain"
	"todo/db-service/domain"
	"todo/pkg/http/handlers"
	"todo/pkg/http/parse"
)

type GetTasksRequest struct{}

type GetTasksResponse struct {
	Tasks []domain.Task `json:"tasks"`
}

type GetTaskByIDRequest struct {
	ID domain.TaskID
}

func CreateGetTaskByIDRequest(r *http.Request) (*GetTaskByIDRequest, error) {
	const queryParamName = "id"

	id, err := parse.Uint64FromQueryParam(r, queryParamName)
	if err != nil {
		return nil, apidomain.ErrInvalidID
	}

	return &GetTaskByIDRequest{ID: id}, nil
}

type GetTaskByIDResponse struct {
	domain.Task
}

type CreateTaskRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func CreateCreateTaskRequest(r *http.Request) (*CreateTaskRequest, error) {
	req := &CreateTaskRequest{}
	if err := handlers.DecodeRequest(r, req); err != nil {
		return nil, apidomain.ErrBadRequest
	}

	if len(req.Name) == 0 {
		return nil, apidomain.ErrEmptyName
	}
	return req, nil
}

type CreateTaskResponse struct {
	ID domain.TaskID `json:"id"`
}

type DeleteTaskByIDRequest struct {
	ID domain.TaskID
}

func CreateDeleteTaskByIDRequest(r *http.Request) (*DeleteTaskByIDRequest, error) {
	const queryParamName = "id"

	id, err := parse.Uint64FromQueryParam(r, queryParamName)
	if err != nil {
		return nil, apidomain.ErrInvalidID
	}

	return &DeleteTaskByIDRequest{ID: id}, nil
}

type DeleteTaskByIDResponse struct{}

type DoneTaskByIDRequest struct {
	ID domain.TaskID
}

func CreateDoneTaskByIDRequest(r *http.Request) (*DoneTaskByIDRequest, error) {
	const queryParamName = "id"

	id, err := parse.Uint64FromQueryParam(r, queryParamName)
	if err != nil {
		return nil, apidomain.ErrInvalidID
	}

	return &DoneTaskByIDRequest{ID: id}, nil
}

type DoneTaskByIDResponse struct{}
