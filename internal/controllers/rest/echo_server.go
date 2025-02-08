package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/wazwki/skillsrock/internal/config"
	"github.com/wazwki/skillsrock/internal/controllers/rest/middlewares"
	"github.com/wazwki/skillsrock/pkg/jwtutil"
)

func NewEchoServer(cfg *config.Config, jwt *jwtutil.JWTUtil) *echo.Echo {
	srv := echo.New()
	srv.HideBanner = true
	srv.GET("/swagger/*", echoSwagger.WrapHandler)
	srv.Use(
		echo.MiddlewareFunc(middlewares.MetricsMiddleware()),
		echo.MiddlewareFunc(middlewares.JWTMiddleware(cfg, jwt)),
		echo.MiddlewareFunc(middlewares.LoggerMiddleware()),
	)

	srv.Server = &http.Server{
		Addr:              fmt.Sprintf("%v:%v", cfg.Host, cfg.Port),
		ReadHeaderTimeout: 800 * time.Millisecond,
		ReadTimeout:       800 * time.Millisecond,
	}

	return srv
}
