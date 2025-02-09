package app

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wazwki/skillsrock/db/postgres"
	"github.com/wazwki/skillsrock/db/redis"
	"github.com/wazwki/skillsrock/internal/config"
	"github.com/wazwki/skillsrock/internal/controllers/rest"
	"github.com/wazwki/skillsrock/internal/controllers/rest/routes"
	v1 "github.com/wazwki/skillsrock/internal/controllers/rest/v1"
	"github.com/wazwki/skillsrock/internal/repository"
	"github.com/wazwki/skillsrock/internal/service"
	"github.com/wazwki/skillsrock/pkg/jwtutil"
	"github.com/wazwki/skillsrock/pkg/logger"
	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
)

type App struct {
	server     *echo.Echo
	migrateDSN string
	pool       *pgxpool.Pool
}

func New(cfg *config.Config) (*App, error) {
	logger.LogInit(cfg.LogLevel)
	logger.Info("App started", zap.String("module", "skillsrock"))

	pool, err := postgres.ConnectPool(cfg.DBdsn)
	if err != nil {
		logger.Error("Fail connect pool", zap.Error(err), zap.String("module", "skillsrock"))
		return nil, err
	}

	redisClient, err := redis.Config(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword, cfg.RedisDBNumber)
	if err != nil {
		logger.Error("Fail connect to redis", zap.Error(err), zap.String("module", "skillsrock"))
		return nil, err
	}

	taskRepository := repository.NewTaskRepository(pool, redisClient)
	taskService := service.NewTaskService(taskRepository)
	taskControllers := v1.NewTaskControllers(taskService)

	userRepository := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepository)
	userControllers := v1.NewUserControllers(userService)

	jwt := jwtutil.NewJWTUtil(jwtutil.Config{
		AccessTokenSecret:  []byte(cfg.AccessTokenSecret),
		RefreshTokenSecret: []byte(cfg.RefreshTokenSecret),
		AccessTokenTTL:     time.Duration(cfg.AccessTokenTTL) * time.Second,
		RefreshTokenTTL:    time.Duration(cfg.RefreshTokenTTL) * time.Second,
	})

	srv := rest.NewEchoServer(cfg, jwt)
	routes.RegisterRoutes(srv, taskControllers, userControllers)

	return &App{server: srv, migrateDSN: cfg.DBdsn, pool: pool}, nil
}

func (a *App) Run() error {
	if err := postgres.RunMigrations(a.migrateDSN); err != nil {
		logger.Error("Fail migrate", zap.Error(err), zap.String("module", "skillsrock"))
		return err
	}

	go func() {
		a.server.Logger.Fatal(a.server.StartServer(a.server.Server))
	}()

	return nil
}

func (a *App) Stop() error {
	logger.Info("App stopping", zap.String("module", "skillsrock"))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.pool.Close()

	if err := a.server.Shutdown(ctx); err != nil {
		logger.Error("Fail to shutdown server", zap.Error(err), zap.String("module", "skillsrock"))
		return err
	}

	logger.Info("App stopped", zap.String("module", "skillsrock"))

	return nil
}
