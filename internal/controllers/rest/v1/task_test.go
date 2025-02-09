package v1_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	v1 "github.com/wazwki/skillsrock/internal/controllers/rest/v1"
	"github.com/wazwki/skillsrock/internal/domain"
	"github.com/wazwki/skillsrock/internal/service/mocks"
)

func TestGetTasks(t *testing.T) {
	mockService := mocks.NewTaskServiceInterface(t)
	server := v1.NewTaskControllers(mockService)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService.On("GetTasks", mock.Anything, mock.Anything).Return([]*domain.Task{}, nil)

	if assert.NoError(t, server.GetTasks(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestCreateTask(t *testing.T) {
	mockService := mocks.NewTaskServiceInterface(t)
	server := v1.NewTaskControllers(mockService)
	e := echo.New()

	taskReq := domain.TaskRequest{Title: "Task 1", Status: "pending", Due_date: "2025-12-12 15:04:05", Priority: "low", Description: "Description 1"}
	jsonReq, _ := json.Marshal(taskReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewReader(jsonReq))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService.On("CreateTask", mock.Anything, mock.Anything).Return("", nil)

	if assert.NoError(t, server.CreateTask(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestUpdateTask(t *testing.T) {
	mockService := mocks.NewTaskServiceInterface(t)
	server := v1.NewTaskControllers(mockService)
	e := echo.New()

	taskReq := domain.TaskRequest{Title: "Task 1", Status: "pending", Due_date: "2025-12-12 15:04:05", Priority: "low", Description: "Description 1"}
	jsonReq, _ := json.Marshal(taskReq)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/tasks/1", bytes.NewReader(jsonReq))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockService.On("UpdateTask", mock.Anything, mock.Anything).Return(&domain.Task{}, nil)

	if assert.NoError(t, server.UpdateTask(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var task domain.TaskResponse
		assert.NoError(t, json.NewDecoder(rec.Body).Decode(&task))
	}
}

// TODO: Test ImportTasks

func TestDeleteTask(t *testing.T) {
	mockService := mocks.NewTaskServiceInterface(t)
	server := v1.NewTaskControllers(mockService)
	e := echo.New()

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockService.On("DeleteTask", mock.Anything, mock.Anything).Return(nil)

	if assert.NoError(t, server.DeleteTask(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

}

func TestGetAnalytics(t *testing.T) {
	mockService := mocks.NewTaskServiceInterface(t)
	server := v1.NewTaskControllers(mockService)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService.On("GetAnalytics", mock.Anything).Return(&domain.Analyse{}, nil)

	if assert.NoError(t, server.GetAnalytics(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

}

func TestImportTasks(t *testing.T) {
	mockService := mocks.NewTaskServiceInterface(t)
	server := v1.NewTaskControllers(mockService)
	e := echo.New()

	tasks := []domain.TaskRequest{
		{Title: "Task 1", Status: "pending", Due_date: "2025-12-12 15:04:05", Priority: "low", Description: "Description 1"},
		{Title: "Task 2", Status: "pending", Due_date: "2025-12-12 15:04:05", Priority: "low", Description: "Description 1"},
	}

	jsonReq, _ := json.Marshal(tasks)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks/import", bytes.NewReader(jsonReq))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService.On("ImportTasks", mock.Anything, mock.Anything).Return(nil)

	if assert.NoError(t, server.ImportTasks(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

}

func TestExportTasks(t *testing.T) {
	mockService := mocks.NewTaskServiceInterface(t)
	server := v1.NewTaskControllers(mockService)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/export", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService.On("ExportTasks", mock.Anything).Return([]*domain.Task{}, nil)

	if assert.NoError(t, server.ExportTasks(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
