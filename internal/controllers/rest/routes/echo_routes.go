package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/wazwki/skillsrock/internal/controllers/rest"
)

func RegisterRoutes(e *echo.Echo, taskControllers rest.TaskControllersInterface, userControllers rest.UserControllersInterface) {
	api := e.Group("/api")
	v1 := api.Group("/v1")

	v1.POST("/auth/register", userControllers.Register)
	v1.POST("/auth/login", userControllers.Login)

	v1.GET("/tasks", taskControllers.GetTasks)
	v1.POST("/tasks", taskControllers.CreateTask)
	v1.PUT("/tasks/:id", taskControllers.UpdateTask)
	v1.DELETE("/tasks/:id", taskControllers.DeleteTask)
	v1.GET("/analytics", taskControllers.GetAnalytics)
	v1.POST("/tasks/import", taskControllers.ImportTasks)
	v1.GET("/tasks/export", taskControllers.ExportTasks)
}
