package domain

import "time"

type Task struct {
	ID          int
	Title       string
	Description string
	Status      string
	Priority    string
	Due_date    time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TaskResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	Due_date    string `json:"due_date"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type TaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	Due_date    string `json:"due_date"`
}

func TaskFromTaskRequest(task *TaskRequest) *Task {
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", task.Due_date)
	return &Task{
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		Due_date:    parsedTime,
	}
}

func TaskToTaskResponse(task *Task) *TaskResponse {
	return &TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		Due_date:    task.Due_date.Format("2006-01-02 15:04:05"),
		CreatedAt:   task.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   task.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func TaskFromTaskResponse(task *TaskResponse) *Task {
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", task.Due_date)
	createdAt, _ := time.Parse("2006-01-02 15:04:05", task.CreatedAt)
	updatedAt, _ := time.Parse("2006-01-02 15:04:05", task.UpdatedAt)
	return &Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		Due_date:    parsedTime,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

type TaskFilter struct {
	Status   string
	SortBy   string
	Priority string
	Name     string
}

type Analyse struct {
	Done        int     `json:"done"`
	InProgress  int     `json:"in_progress"`
	Pending     int     `json:"pending"`
	AverageTime float64 `json:"average_time"`
	Weekly      WeeklyReport
}

type WeeklyReport struct {
	Completed   int `json:"completed"`
	Uncompleted int `json:"uncompleted"`
}
