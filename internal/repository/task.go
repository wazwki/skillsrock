package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/wazwki/skillsrock/internal/domain"
)

type TaskRepository struct {
	DataBase *pgxpool.Pool
	Cache    *redis.Client
}

func NewTaskRepository(db *pgxpool.Pool, cache *redis.Client) TaskRepositoryInterface {
	return &TaskRepository{DataBase: db, Cache: cache}
}

func (r *TaskRepository) CreateTask(ctx context.Context, task *domain.Task) (string, error) {
	query := `INSERT INTO tasks (title, description, status, priority, due_date) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id string

	err := r.DataBase.QueryRow(ctx, query, task.Title, task.Description, task.Status, task.Priority, task.Due_date).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *TaskRepository) GetTasks(ctx context.Context, filter domain.TaskFilter) ([]*domain.Task, error) {
	query := `SELECT id, title, description, status, priority, due_date FROM tasks WHERE 1=1`
	args := make([]interface{}, 0)

	if filter.Status != "" {
		query += " AND status = $1"
		args = append(args, filter.Status)
	}

	if filter.Priority != "" {
		query += " AND priority = $2"
		args = append(args, filter.Priority)
	}

	if filter.Name != "" {
		query += " AND title = $3"
		args = append(args, filter.Name)
	}

	if filter.SortBy != "" {
		switch filter.SortBy {
		case "low":
			query += " ORDER BY due_date ASC"
		case "max":
			query += " ORDER BY due_date DESC"
		}
	}

	rows, err := r.DataBase.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := make([]*domain.Task, 0)
	for rows.Next() {
		task := &domain.Task{}
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.Due_date)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskRepository) UpdateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	query := `UPDATE tasks SET title = $1, description = $2, status = $3, priority = $4, due_date = $5 WHERE id = $6 RETURNING id, title, description, status, priority, due_date`
	err := r.DataBase.QueryRow(ctx, query, task.Title, task.Description, task.Status, task.Priority, task.Due_date, task.ID).Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.Due_date)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *TaskRepository) DeleteTask(ctx context.Context, task_id string) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.DataBase.Exec(ctx, query, task_id)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepository) ClearTasks(ctx context.Context) error {
	query := `DELETE FROM tasks WHERE status != 'done' AND due_date < CURRENT_DATE - INTERVAL '7 days'`
	if _, err := r.DataBase.Exec(ctx, query); err != nil {
		return err
	}
	return nil
}

func (r *TaskRepository) GetCachedAnalytics(ctx context.Context) (*domain.Analyse, error) {
	val, err := r.Cache.Get(ctx, "analytics").Result()
	if err != nil {
		return nil, err
	}

	task := &domain.Analyse{}
	err = json.Unmarshal([]byte(val), task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *TaskRepository) GetAnalytics(ctx context.Context) (*domain.Analyse, error) {
	var analyse domain.Analyse
	var week domain.WeeklyReport

	query := `SELECT COUNT(*) FROM tasks WHERE status = 'done' AND due_date >= CURRENT_DATE - INTERVAL '7 days'`
	err := r.DataBase.QueryRow(ctx, query).Scan(&week.Completed)
	if err != nil {
		return nil, err
	}

	query = `SELECT COUNT(*) FROM tasks WHERE status != 'done' AND due_date >= CURRENT_DATE - INTERVAL '7 days'`
	err = r.DataBase.QueryRow(ctx, query).Scan(&week.Uncompleted)
	if err != nil {
		return nil, err
	}

	analyse.Weekly = week

	query = `SELECT status, COUNT(*) FROM tasks GROUP BY status`
	rows, err := r.DataBase.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		err := rows.Scan(&status, &count)
		if err != nil {
			return nil, err
		}

		switch status {
		case "done":
			analyse.Done = count
		case "in_progress":
			analyse.InProgress = count
		case "pending":
			analyse.Pending = count
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	query = `SELECT AVG(created_at - due_date) FROM tasks WHERE status = 'done'`
	err = r.DataBase.QueryRow(ctx, query).Scan(&analyse.AverageTime)
	if err != nil {
		return nil, err
	}

	return &analyse, nil
}

func (r *TaskRepository) SetAnalytics(ctx context.Context, task *domain.Analyse) error {

	data, err := json.Marshal(*task)
	if err != nil {
		return err
	}

	err = r.Cache.Set(ctx, "analytics", data, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepository) ImportTasks(ctx context.Context, task []*domain.Task) error {
	tx, err := r.DataBase.Begin(ctx)
	if err != nil {
		return err
	}

	for _, t := range task {
		_, err := tx.Exec(ctx, `INSERT INTO tasks (title, description, status, priority, due_date) VALUES ($1, $2, $3, $4, $5)`, t.Title, t.Description, t.Status, t.Priority, t.Due_date)
		if err != nil {
			err = tx.Rollback(ctx)
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) ExportTasks(ctx context.Context) ([]*domain.Task, error) {
	query := `SELECT id, title, description, status, priority, due_date FROM tasks`
	rows, err := r.DataBase.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := make([]*domain.Task, 0)
	for rows.Next() {
		task := &domain.Task{}
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.Due_date)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
