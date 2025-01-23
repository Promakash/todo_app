package domain

import (
	"fmt"
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
