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

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) CheckUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func TestRegisterUser(t *testing.T) {
	e := echo.New()
	mockService := new(MockUserService)
	controller := v1.NewUserControllers(mockService)
	user := &domain.User{Name: "John Doe", Password: "password123"}
	mockService.On("CreateUser", mock.Anything, user).Return(user, nil)

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, controller.Register(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestLoginUser(t *testing.T) {
	e := echo.New()
	mockService := new(MockUserService)
	controller := v1.NewUserControllers(mockService)
	user := &domain.User{Name: "John Doe", Password: "password123"}
	mockService.On("CheckUser", mock.Anything, user).Return(nil)

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, controller.Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
