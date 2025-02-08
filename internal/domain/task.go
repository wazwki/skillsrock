package domain

type Task struct {
	ID          int
	Title       string
	Description string
	Status      string
	Priority    string
	Due_date    string
	CreatedAt   string
	UpdatedAt   string
}

type TaskDTO struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	Due_date    string `json:"due_date"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func TaskFromTaskDTO(taskDTO *TaskDTO) *Task {
	return &Task{
		ID:          taskDTO.ID,
		Title:       taskDTO.Title,
		Description: taskDTO.Description,
		Status:      taskDTO.Status,
		Priority:    taskDTO.Priority,
		Due_date:    taskDTO.Due_date,
		CreatedAt:   taskDTO.CreatedAt,
		UpdatedAt:   taskDTO.UpdatedAt,
	}
}

func TaskDTOFromTask(task *Task) *TaskDTO {
	return &TaskDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		Due_date:    task.Due_date,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

type TaskFilter struct {
	Status   string
	SortBy   string
	Priority string
	Name     string
}

type Analyse struct {
	Done        int `json:"done"`
	InProgress  int `json:"in_progress"`
	Pending     int `json:"pending"`
	AverageTime int `json:"average_time"`
	Weekly      WeeklyReport
}

type WeeklyReport struct {
	Completed   int `json:"completed"`
	Uncompleted int `json:"uncompleted"`
}
