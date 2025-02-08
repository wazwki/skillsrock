package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wazwki/skillsrock/internal/domain"
)

type UserRepository struct {
	DataBase *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepositoryInterface {
	return &UserRepository{DataBase: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `INSERT INTO users (name, password) VALUES ($1, $2) RETURNING id`

	err := r.DataBase.QueryRow(ctx, query, user.Name, user.Password).Scan(&user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) CheckUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `SELECT id, name, password FROM users WHERE name = $1`

	var dbUser domain.User
	err := r.DataBase.QueryRow(ctx, query, user.Name).Scan(&dbUser.ID, &dbUser.Name, &dbUser.Password)
	if err != nil {
		return nil, err
	}

	return &dbUser, nil
}
