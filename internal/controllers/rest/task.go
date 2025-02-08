package rest

import (
	"github.com/labstack/echo/v4"
)

type TaskControllersInterface interface {
	GetTasks(c echo.Context) error
	CreateTask(c echo.Context) error
	UpdateTask(c echo.Context) error
	DeleteTask(c echo.Context) error
	GetAnalytics(c echo.Context) error
	ImportTasks(c echo.Context) error
	ExportTasks(c echo.Context) error
}
