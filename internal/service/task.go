package service

import (
	"context"
	"time"

	"github.com/wazwki/skillsrock/internal/domain"
	"github.com/wazwki/skillsrock/internal/repository"
	"github.com/wazwki/skillsrock/pkg/logger"
	"go.uber.org/zap"
)

type TaskService struct {
	repo repository.TaskRepositoryInterface
}

func NewTaskService(repo repository.TaskRepositoryInterface) TaskServiceInterface {
	t := &TaskService{repo: repo}

	go t.analyseWorker(time.Hour*6, 3, time.Second*5)
	go t.updateWorker(time.Hour*24, 3, time.Second*5)

	return t
}

func (s *TaskService) CreateTask(ctx context.Context, task *domain.Task) (string, error) {
	id, err := s.repo.CreateTask(ctx, task)
	if err != nil {
		logger.Error("Failed to create task", zap.Error(err), zap.String("module", "skillsrock"))
		return "", err
	}

	return id, nil
}

func (s *TaskService) GetTasks(ctx context.Context, filter domain.TaskFilter) ([]*domain.Task, error) {
	tasks, err := s.repo.GetTasks(ctx, filter)
	if err != nil {
		logger.Error("Failed to get tasks", zap.Error(err), zap.String("module", "skillsrock"))
		return nil, err
	}

	return tasks, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	updatedTask, err := s.repo.UpdateTask(ctx, task)
	if err != nil {
		logger.Error("Failed to update task", zap.Error(err), zap.String("module", "skillsrock"))
		return nil, err
	}

	return updatedTask, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, task_id string) error {
	err := s.repo.DeleteTask(ctx, task_id)
	if err != nil {
		logger.Error("Failed to delete task", zap.Error(err), zap.String("module", "skillsrock"))
		return err
	}

	return nil
}

func (s *TaskService) GetAnalytics(ctx context.Context) (*domain.Analyse, error) {
	analytics, err := s.repo.GetCachedAnalytics(ctx)
	if err != nil {
		logger.Error("Failed to get analytics", zap.Error(err), zap.String("module", "skillsrock"))
		return nil, err
	}

	return analytics, nil
}

func (s *TaskService) ImportTasks(ctx context.Context, task []*domain.Task) error {
	err := s.repo.ImportTasks(ctx, task)
	if err != nil {
		logger.Error("Failed to import tasks", zap.Error(err), zap.String("module", "skillsrock"))
		return err
	}
	return nil
}

func (s *TaskService) ExportTasks(ctx context.Context) ([]*domain.Task, error) {
	tasks, err := s.repo.ExportTasks(ctx)
	if err != nil {
		logger.Error("Failed to export tasks", zap.Error(err), zap.String("module", "skillsrock"))
		return nil, err
	}
	return tasks, nil
}

func (s *TaskService) analyseWorker(updateInterval time.Duration, retryCount int, retryInterval time.Duration) {
	tick := time.NewTicker(updateInterval)
	defer tick.Stop()
	for range tick.C {
		var analitics *domain.Analyse
		var err error
		for range retryCount {
			analitics, err = s.repo.GetAnalytics(context.Background())
			if err != nil {
				logger.Error("Failed to get analytics", zap.Error(err), zap.String("module", "skillsrock"))
				time.Sleep(retryInterval)
				continue
			} else {
				break
			}
		}

		for range retryCount {
			err = s.repo.SetAnalytics(context.Background(), analitics)
			if err != nil {
				logger.Error("Failed to set analytics", zap.Error(err), zap.String("module", "skillsrock"))
				time.Sleep(retryInterval)
				continue
			} else {
				break
			}
		}
	}
}

func (s *TaskService) updateWorker(updateInterval time.Duration, retryCount int, retryInterval time.Duration) {
	tick := time.NewTicker(updateInterval)
	defer tick.Stop()
	for range tick.C {
		for range retryCount {
			err := s.repo.ClearTasks(context.Background())
			if err != nil {
				logger.Error("Failed to clear tasks", zap.Error(err), zap.String("module", "skillsrock"))
				time.Sleep(retryInterval)
				continue
			} else {
				break
			}
		}
	}
}
