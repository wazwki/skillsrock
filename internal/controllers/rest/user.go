package rest

import (
	"github.com/labstack/echo/v4"
)

type UserControllersInterface interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
}
