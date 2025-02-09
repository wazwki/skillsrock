package v1

import (
	"encoding/json"
	"net/http"

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
// @Param user body domain.UserRequest true "User"
// @Success 201 {object} domain.UserResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/v1/auth/register [post]
func (s *UserServer) Register(c echo.Context) error {
	var user *domain.UserRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	createdUser, err := s.service.CreateUser(c.Request().Context(), domain.UserRequestToUser(user))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create user"})
	}

	return c.JSON(http.StatusCreated, domain.UserToUserResponse(createdUser))
}

// @Summary Login user
// @Description Login user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body domain.UserRequest true "User"
// @Success 200 {object} nil
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/v1/auth/login [post]
func (s *UserServer) Login(c echo.Context) error {
	var user *domain.UserRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	err := s.service.CheckUser(c.Request().Context(), domain.UserRequestToUser(user))
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}

	return c.NoContent(http.StatusOK)
}
