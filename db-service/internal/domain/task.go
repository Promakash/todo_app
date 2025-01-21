package domain

type TaskID = uint64

type Task struct {
	ID          TaskID `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsDone      bool   `json:"is_done"`
}
