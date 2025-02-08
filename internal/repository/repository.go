package repository

import (
	"context"

	"github.com/wazwki/skillsrock/internal/domain"
)

type TaskRepositoryInterface interface {
	CreateTask(ctx context.Context, task *domain.Task) (string, error)
	GetTasks(ctx context.Context, filter domain.TaskFilter) ([]*domain.Task, error)
	UpdateTask(ctx context.Context, task *domain.Task) (*domain.Task, error)
	DeleteTask(ctx context.Context, task_id string) error
	ClearTasks(ctx context.Context) error
	GetCachedAnalytics(ctx context.Context) (*domain.Analyse, error)
	GetAnalytics(ctx context.Context) (*domain.Analyse, error)
	SetAnalytics(ctx context.Context, task *domain.Analyse) error
	ImportTasks(ctx context.Context, task []*domain.Task) error
	ExportTasks(ctx context.Context) ([]*domain.Task, error)
}

type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	CheckUser(ctx context.Context, user *domain.User) (*domain.User, error)
}
