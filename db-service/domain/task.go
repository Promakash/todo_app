package domain

import (
	"fmt"
	todov1 "todo/protos/gen/go"
)

type TaskID = uint64

func TaskIDToString(id TaskID) string {
	return fmt.Sprintf("%d", id)
}

type Task struct {
	ID          TaskID `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsDone      bool   `json:"is_done"`
}

func (t Task) ToGRPC() *todov1.Task {
	return &todov1.Task{
		Id:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		IsDone:      t.IsDone,
	}
}

func TaskFromGRPC(task *todov1.Task) Task {
	return Task{
		ID:          task.Id,
		Name:        task.Name,
		Description: task.Description,
		IsDone:      task.IsDone,
	}
}
