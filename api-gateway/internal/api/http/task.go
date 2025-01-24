package http

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"todo/api-gateway/internal/api/http/types"
	"todo/api-gateway/internal/clients/todo/grpc"
	"todo/api-gateway/internal/domain"
	"todo/pkg/http/handlers"
	"todo/pkg/http/responses"
)

type TaskHandler struct {
	logger  *slog.Logger
	grpcAPI *grpc.Client
}

func NewTaskHandler(logger *slog.Logger, grpcAPI *grpc.Client) *TaskHandler {
	return &TaskHandler{
		logger:  logger,
		grpcAPI: grpcAPI,
	}
}

const singleTaskRequestPath = "/task"
const multipleTasksRequestPath = "/tasks"

func (h *TaskHandler) WithTaskHandlers() handlers.RouterOption {
	return func(r chi.Router) {
		r.Route(
			"/",
			func(r chi.Router) {
				handlers.AddHandler(r.Get, multipleTasksRequestPath, h.getTasksHandler)
				handlers.AddHandler(r.Get, singleTaskRequestPath, h.getTaskByIDHandler)
				handlers.AddHandler(r.Delete, singleTaskRequestPath, h.deleteTaskByIDHandler)
				handlers.AddHandler(r.Post, singleTaskRequestPath, h.CreateTaskHandler)
				handlers.AddHandler(r.Put, singleTaskRequestPath, h.DoneTaskByIDHandler)
			},
		)
	}
}

// @Summary		Returns array of all tasks
// @Description	Returns array of all tasks that exist in DB
// @Produce		json
// @Success		200	{object}	types.GetTasksResponse	"Successfully returned"
// @Failure		500	{object}	responses.ErrorResponse	"Internal server error"
// @Router			/tasks [get]
func (h *TaskHandler) getTasksHandler(r *http.Request) responses.Response {
	tasks, err := h.grpcAPI.ListTasks(r.Context())
	return domain.HandleError(err, &types.GetTasksResponse{Tasks: tasks})
}

// @Summary		Returns task by ID
// @Description	Returns task by ID
// @Produce		json
// @Param			id	query		integer						true	"Task ID"	format(int64)	minimum(1)
// @Success		200	{object}	types.GetTaskByIDResponse	"Task returned successfully"
// @Failure		400	{object}	responses.ErrorResponse		"Bad Request"
// @Failure		404	{object}	responses.ErrorResponse		"Task does not exist"
// @Failure		500	{object}	responses.ErrorResponse		"Internal server error
// @Router			/task [get]
func (h *TaskHandler) getTaskByIDHandler(r *http.Request) responses.Response {
	req, err := types.CreateGetTaskByIDRequest(r)
	if err != nil {
		return domain.HandleError(err, nil)
	}

	task, err := h.grpcAPI.GetTaskByID(r.Context(), req.ID)
	return domain.HandleError(err, &types.GetTaskByIDResponse{Task: task})
}

// @Summary		Deletes task by ID
// @Description	Deletes task by ID
// @Produce		json
// @Param			id	query		integer							true	"Task ID"	format(int64)	minimum(1)
// @Success		200	{object}	types.DeleteTaskByIDResponse	"Task deleted successfully"
// @Failure		400	{object}	responses.ErrorResponse			"Bad Request"
// @Failure		404	{object}	responses.ErrorResponse			"Task does not exist"
// @Failure		500	{object}	responses.ErrorResponse			"Internal server error
// @Router			/task [delete]
func (h *TaskHandler) deleteTaskByIDHandler(r *http.Request) responses.Response {
	req, err := types.CreateDeleteTaskByIDRequest(r)
	if err != nil {
		return domain.HandleError(err, nil)
	}

	err = h.grpcAPI.DeleteTaskByID(r.Context(), req.ID)
	return domain.HandleError(err, nil)
}

// @Summary		Creates task
// @Description	Creates task
// @Produce		json
// @Accept			json
// @Param			request	body		types.CreateTaskRequest		true	"Task info"
// @Success		200		{object}	types.CreateTaskResponse	"Task created successfully"
// @Failure		400		{object}	responses.ErrorResponse		"Bad Request"
// @Failure		500		{object}	responses.ErrorResponse		"Internal server error
// @Router			/task [post]
func (h *TaskHandler) CreateTaskHandler(r *http.Request) responses.Response {
	req, err := types.CreateCreateTaskRequest(r)
	if err != nil {
		return domain.HandleError(err, nil)
	}

	id, err := h.grpcAPI.CreateTask(r.Context(), req.Name, req.Description)
	return domain.HandleError(err, &types.CreateTaskResponse{ID: id})
}

// @Summary		Marks task as done
// @Description	Marks task as done
// @Produce		json
// @Param			id	query		integer						true	"Task ID"	format(int64)	minimum(1)
// @Success		200	{object}	types.DoneTaskByIDResponse	"Task marked as done successfully"
// @Failure		400	{object}	responses.ErrorResponse		"Bad Request"
// @Failure		404	{object}	responses.ErrorResponse		"Task does not exist"
// @Failure		500	{object}	responses.ErrorResponse		"Internal server error
// @Router			/task [put]
func (h *TaskHandler) DoneTaskByIDHandler(r *http.Request) responses.Response {
	req, err := types.CreateDoneTaskByIDRequest(r)
	if err != nil {
		return domain.HandleError(err, nil)
	}

	err = h.grpcAPI.DoneTaskByID(r.Context(), req.ID)
	return domain.HandleError(err, nil)
}
