package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/wazwki/skillsrock/internal/controllers/rest"
	"github.com/wazwki/skillsrock/internal/domain"
	"github.com/wazwki/skillsrock/internal/service"
)

type TaskServer struct {
	service service.TaskServiceInterface
}

func NewTaskControllers(s service.TaskServiceInterface) rest.TaskControllersInterface {
	return &TaskServer{service: s}
}

// @Summary Get tasks
// @Description Get tasks
// @Tags Tasks
// @Accept json
// @Produce json
// @Param status query string false "Choose status: pending, in_progress, done"
// @Param sort_by query string false "Choose sort by date: low, high"
// @Param priority query string false "Choose priority: low, medium, high"
// @Param name query string false "Choose name"
// @Success 200 {object} []domain.TaskResponse
// @Failure 500 {object} string
// @Router /api/v1/tasks [get]
func (s *TaskServer) GetTasks(c echo.Context) error {
	status := c.QueryParam("status")
	sortBy := c.QueryParam("sort_by")
	priority := c.QueryParam("priority")
	name := c.QueryParam("name")

	tasks, err := s.service.GetTasks(c.Request().Context(), domain.TaskFilter{
		Status:   status,
		SortBy:   sortBy,
		Priority: priority,
		Name:     name,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to get tasks"})
	}

	var tasksR []*domain.TaskResponse
	for _, task := range tasks {
		tasksR = append(tasksR, domain.TaskToTaskResponse(task))
	}

	return c.JSON(http.StatusOK, tasksR)
}

// @Summary Create task
// @Description Create task
// @Tags Tasks
// @Accept json
// @Produce json
// @Param task body domain.TaskRequest true "Task"
// @Success 201 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/v1/tasks [post]
func (s *TaskServer) CreateTask(c echo.Context) error {
	var task *domain.TaskRequest

	err := json.NewDecoder(c.Request().Body).Decode(&task)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	id, err := s.service.CreateTask(c.Request().Context(), domain.TaskFromTaskRequest(task))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create task"})
	}

	return c.JSON(http.StatusCreated, id)
}

// @Summary Update task
// @Description Update task
// @Tags Tasks
// @Accept json
// @Produce json
// @Param task body domain.TaskRequest true "Task"
// @Param id path string true "ID"
// @Success 200 {object} domain.TaskResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/v1/tasks/{id} [put]
func (s *TaskServer) UpdateTask(c echo.Context) error {
	var task *domain.TaskRequest

	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	err = json.NewDecoder(c.Request().Body).Decode(&task)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	dTask := domain.TaskFromTaskRequest(task)
	dTask.ID = id

	uTask, err := s.service.UpdateTask(c.Request().Context(), dTask)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update task"})
	}

	return c.JSON(http.StatusOK, domain.TaskToTaskResponse(uTask))
}

// @Summary Delete task
// @Description Delete task
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} nil
// @Failure 500 {object} string
// @Router /api/v1/tasks/{id} [delete]
func (s *TaskServer) DeleteTask(c echo.Context) error {
	id := c.Param("id")

	err := s.service.DeleteTask(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete task"})
	}

	return nil
}

// @Summary Get analytics
// @Description Get analytics
// @Tags Analytics
// @Accept json
// @Produce json
// @Success 200 {object} domain.Analyse
// @Failure 500 {object} string
// @Router /api/v1/analytics [get]
func (s *TaskServer) GetAnalytics(c echo.Context) error {
	var analytics *domain.Analyse
	analytics, err := s.service.GetAnalytics(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to get analytics"})
	}

	return c.JSON(http.StatusOK, analytics)
}

// @Summary Import tasks
// @Description Import tasks
// @Tags Tasks
// @Accept json
// @Produce json
// @Param tasks body []domain.TaskRequest true "Tasks"
// @Success 200 {object} nil
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/v1/tasks/import [post]
func (s *TaskServer) ImportTasks(c echo.Context) error {
	var tasks []*domain.TaskRequest

	err := json.NewDecoder(c.Request().Body).Decode(&tasks)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	if len(tasks) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "No tasks provided"})
	}

	tasksD := make([]*domain.Task, 0, len(tasks))

	for _, task := range tasks {
		if task == nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
		}

		taskD := domain.TaskFromTaskRequest(task)
		if taskD == nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid task data"})
		}

		tasksD = append(tasksD, taskD)
	}

	err = s.service.ImportTasks(c.Request().Context(), tasksD)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to import tasks"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Tasks imported successfully"})
}

// @Summary Export tasks
// @Description Export tasks
// @Tags Tasks
// @Accept json
// @Produce json
// @Success 200 {object} []domain.TaskResponse
// @Failure 500 {object} string
// @Router /api/v1/tasks/export [get]
func (s *TaskServer) ExportTasks(c echo.Context) error { // Service
	tasks, err := s.service.ExportTasks(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to export tasks"})
	}

	tasksR := make([]*domain.TaskResponse, len(tasks))
	for _, task := range tasks {
		tasksR = append(tasksR, domain.TaskToTaskResponse(task))
	}

	return c.JSON(http.StatusOK, tasksR)
}
