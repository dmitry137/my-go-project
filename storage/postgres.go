package storage

import (
	"context"
	"errors"

	"github.com/dmitry137/my-go-project/models"
	"github.com/jackc/pgx/v4"
)

type Storage interface {
	CreateTask(task models.Task) (models.Task, error)
	GetTasks() ([]models.Task, error)
	GetTaskByID(id int) (models.Task, error)
	UpdateTask(task models.Task) error
	DeleteTask(id int) error
}

type PostgresStorage struct {
	conn *pgx.Conn
}

func NewPostgresStorage(conn *pgx.Conn) *PostgresStorage {
	return &PostgresStorage{conn: conn}
}

func (s *PostgresStorage) CreateTask(task models.Task) (models.Task, error) {
	query := `INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := s.conn.QueryRow(context.Background(), query, task.Title, task.Description, task.Status).Scan(
		&task.ID, &task.CreatedAt, &task.UpdatedAt,
	)
	return task, err
}

func (s *PostgresStorage) GetTasks() ([]models.Task, error) {
	query := `SELECT id, title, description, status, created_at, updated_at FROM tasks`
	rows, err := s.conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var description *string
		err := rows.Scan(
			&task.ID, &task.Title, &description, &task.Status, &task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		task.Description = description
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *PostgresStorage) GetTaskByID(id int) (models.Task, error) {
	query := `SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE id = $1`
	var task models.Task
	var description *string
	err := s.conn.QueryRow(context.Background(), query, id).Scan(
		&task.ID, &task.Title, &description, &task.Status, &task.CreatedAt, &task.UpdatedAt,
	)
	task.Description = description
	if errors.Is(err, pgx.ErrNoRows) {
		return task, errors.New("task not found")
	}
	return task, err
}

func (s *PostgresStorage) UpdateTask(task models.Task) error {
	query := `UPDATE tasks SET title=$1, description=$2, status=$3, updated_at=NOW() WHERE id=$4`
	_, err := s.conn.Exec(context.Background(), query, task.Title, task.Description, task.Status, task.ID)
	return err
}

func (s *PostgresStorage) DeleteTask(id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	res, err := s.conn.Exec(context.Background(), query, id)
	if res.RowsAffected() == 0 {
		return errors.New("task not found")
	}
	return err
}
