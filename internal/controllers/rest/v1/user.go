package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/wazwki/skillsrock/internal/controllers/rest"
	"github.com/wazwki/skillsrock/internal/domain"
	"github.com/wazwki/skillsrock/internal/service"
)

type UserServer struct {
	service service.UserServiceInterface
}

func NewUserControllers(s service.UserServiceInterface) rest.UserControllersInterface {
	return &UserServer{service: s}
}

// @Summary Register user
// @Description Register user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body domain.User true "User"
// @Success 201 {object} domain.User
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/v1/users [post]
func (s *UserServer) Register(c echo.Context) error {
	var user *domain.User
	if err := json.NewDecoder(c.Request().Body).Decode(user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	createdUser, err := s.service.CreateUser(c.Request().Context(), user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create user"})
	}

	return c.JSON(http.StatusCreated, createdUser)
}

// @Summary Login user
// @Description Login user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body domain.User true "User"
// @Success 200 {object} nil
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/v1/users/login [post]
func (s *UserServer) Login(c echo.Context) error {
	var user *domain.User
	if err := json.NewDecoder(c.Request().Body).Decode(user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	err := s.service.CheckUser(c.Request().Context(), user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "User not found"})
	}

	c.SetCookie(&http.Cookie{
		Name:     "Authorization",
		Value:    c.Get("Authorization").(string),
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return c.NoContent(http.StatusOK)
}
