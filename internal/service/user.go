package service

import (
	"context"

	"github.com/wazwki/skillsrock/internal/domain"
	"github.com/wazwki/skillsrock/internal/repository"
	"github.com/wazwki/skillsrock/pkg/hashutil"
	"github.com/wazwki/skillsrock/pkg/logger"
	"go.uber.org/zap"
)

type UserService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) UserServiceInterface {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	hashed, err := hashutil.HashPassword(user.Password)
	if err != nil {
		logger.Error("Failed to hash password", zap.Error(err), zap.String("module", "skillsrock"))
		return nil, err
	}
	user.Password = hashed

	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) CheckUser(ctx context.Context, user *domain.User) error {
	dbUser, err := s.repo.CheckUser(ctx, user)
	if err != nil {
		logger.Error("Failed to check user", zap.Error(err), zap.String("module", "skillsrock"))
		return err
	}

	if hashutil.ComparePassword(dbUser.Password, user.Password) {
		return nil
	}

	return domain.UserNotFound
}
