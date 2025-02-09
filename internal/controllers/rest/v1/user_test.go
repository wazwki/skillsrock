package v1_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	v1 "github.com/wazwki/skillsrock/internal/controllers/rest/v1"
	"github.com/wazwki/skillsrock/internal/domain"
	"github.com/wazwki/skillsrock/internal/service/mocks"
)

func TestRegisterUser(t *testing.T) {
	mockService := mocks.NewUserServiceInterface(t)
	server := v1.NewUserControllers(mockService)
	e := echo.New()

	userReq := domain.UserRequest{Name: "John Doe", Password: "password"}
	jsonReq, _ := json.Marshal(userReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(jsonReq))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService.On("CreateUser", mock.Anything, mock.Anything).Return(&domain.User{}, nil)

	if assert.NoError(t, server.Register(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestLoginUser(t *testing.T) {
	mockService := mocks.NewUserServiceInterface(t)
	server := v1.NewUserControllers(mockService)
	e := echo.New()

	userReq := domain.UserRequest{Name: "John Doe", Password: "password"}
	jsonReq, _ := json.Marshal(userReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(jsonReq))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService.On("CheckUser", mock.Anything, mock.Anything).Return(nil)

	if assert.NoError(t, server.Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
