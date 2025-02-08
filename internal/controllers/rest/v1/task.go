package v1

import (
	"encoding/json"
	"net/http"

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
// @Param status query string false "Status"
// @Param sort_by query string false "Sort by"
// @Param priority query string false "Priority"
// @Param name query string false "Name"
// @Success 200 {object} []domain.Task
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

	return c.JSON(http.StatusOK, tasks)
}

// @Summary Create task
// @Description Create task
// @Tags Tasks
// @Accept json
// @Produce json
// @Param task body domain.Task true "Task"
// @Success 201 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/v1/tasks [post]
func (s *TaskServer) CreateTask(c echo.Context) error {
	var task *domain.Task

	err := json.NewDecoder(c.Request().Body).Decode(task)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	id, err := s.service.CreateTask(c.Request().Context(), task)
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
// @Param task body domain.Task true "Task"
// @Success 200 {object} domain.Task
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/v1/tasks/{id} [put]
func (s *TaskServer) UpdateTask(c echo.Context) error {
	var task *domain.Task

	err := json.NewDecoder(c.Request().Body).Decode(task)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	uTask, err := s.service.UpdateTask(c.Request().Context(), task)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update task"})
	}

	return c.JSON(http.StatusOK, uTask)
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
// @Param tasks body []domain.Task true "Tasks"
// @Success 200 {object} nil
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/v1/tasks/import [post]
func (s *TaskServer) ImportTasks(c echo.Context) error {
	var tasks []*domain.Task

	err := json.NewDecoder(c.Request().Body).Decode(&tasks)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	err = s.service.ImportTasks(c.Request().Context(), tasks)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to import tasks"})
	}

	return nil
}

// @Summary Export tasks
// @Description Export tasks
// @Tags Tasks
// @Accept json
// @Produce json
// @Success 200 {object} []domain.Task
// @Failure 500 {object} string
// @Router /api/v1/tasks/export [get]
func (s *TaskServer) ExportTasks(c echo.Context) error {
	tasks, err := s.service.ExportTasks(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to export tasks"})
	}

	return c.JSON(http.StatusOK, tasks)
}
