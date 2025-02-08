package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	v1 "github.com/wazwki/skillsrock/internal/controllers/rest/v1"
	"github.com/wazwki/skillsrock/internal/domain"
)

type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) GetTasks(ctx context.Context, filter domain.TaskFilter) ([]*domain.Task, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*domain.Task), args.Error(1)
}

func (m *MockTaskService) CreateTask(ctx context.Context, task *domain.Task) (string, error) {
	args := m.Called(ctx, task)
	return args.String(0), args.Error(1)
}

func (m *MockTaskService) UpdateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	args := m.Called(ctx, task)
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskService) DeleteTask(ctx context.Context, task_id string) error {
	args := m.Called(ctx, task_id)
	return args.Error(0)
}

func (m *MockTaskService) GetAnalytics(ctx context.Context) (*domain.Analyse, error) {
	args := m.Called(ctx)
	return args.Get(0).(*domain.Analyse), args.Error(1)
}

func (m *MockTaskService) ImportTasks(ctx context.Context, tasks []*domain.Task) error {
	args := m.Called(ctx, tasks)
	return args.Error(0)
}

func (m *MockTaskService) ExportTasks(ctx context.Context) ([]*domain.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Task), args.Error(1)
}

func TestGetTasks(t *testing.T) {
	e := echo.New()
	mockService := new(MockTaskService)
	controller := v1.NewTaskControllers(mockService)
	tasks := []*domain.Task{{ID: 1, Title: "Test Task", Description: "Test Description", Status: "Pending", Priority: "High", Due_date: "2025-12-01", CreatedAt: "2025-01-01", UpdatedAt: "2025-01-02"}}
	mockService.On("GetTasks", mock.Anything, mock.Anything).Return(tasks, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, controller.GetTasks(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestCreateTask(t *testing.T) {
	e := echo.New()
	mockService := new(MockTaskService)
	controller := v1.NewTaskControllers(mockService)
	task := &domain.Task{ID: 2, Title: "New Task", Description: "New Description", Status: "InProgress", Priority: "Medium", Due_date: "2025-12-10", CreatedAt: "2025-02-01", UpdatedAt: "2025-02-02"}
	mockService.On("CreateTask", mock.Anything, task).Return("123", nil)

	body, _ := json.Marshal(task)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, controller.CreateTask(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestUpdateTask(t *testing.T) {
	e := echo.New()
	mockService := new(MockTaskService)
	controller := v1.NewTaskControllers(mockService)
	task := &domain.Task{ID: 2, Title: "Updated Task", Description: "Updated Description", Status: "Done", Priority: "Low", Due_date: "2025-12-10", CreatedAt: "2025-02-01", UpdatedAt: "2025-02-02"}
	mockService.On("UpdateTask", mock.Anything, task).Return(task, nil)

	body, _ := json.Marshal(task)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/tasks/2", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("2")

	if assert.NoError(t, controller.UpdateTask(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestDeleteTask(t *testing.T) {
	e := echo.New()
	mockService := new(MockTaskService)
	controller := v1.NewTaskControllers(mockService)
	mockService.On("DeleteTask", mock.Anything, "1").Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, controller.DeleteTask(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestGetAnalytics(t *testing.T) {
	e := echo.New()
	mockService := new(MockTaskService)
	controller := v1.NewTaskControllers(mockService)
	analytics := &domain.Analyse{Done: 5, InProgress: 3, Pending: 2, AverageTime: 24, Weekly: domain.WeeklyReport{Completed: 10, Uncompleted: 5}}
	mockService.On("GetAnalytics", mock.Anything).Return(analytics, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, controller.GetAnalytics(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestImportTasks(t *testing.T) {
	e := echo.New()
	mockService := new(MockTaskService)
	controller := v1.NewTaskControllers(mockService)
	mockService.On("ImportTasks", mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks/import", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, controller.ImportTasks(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestExportTasks(t *testing.T) {
	e := echo.New()
	mockService := new(MockTaskService)
	controller := v1.NewTaskControllers(mockService)
	tasks := []*domain.Task{{ID: 1, Title: "Exported Task", Description: "Exported Description", Status: "Done", Priority: "Low", Due_date: "2025-11-15", CreatedAt: "2025-03-01", UpdatedAt: "2025-03-02"}}
	mockService.On("ExportTasks", mock.Anything).Return(tasks, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/export", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, controller.ExportTasks(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
