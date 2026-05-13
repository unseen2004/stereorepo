package tasks

import (
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
)

type Task struct {
	ID              string           `db:"id" json:"id"`
	Title           string           `db:"title" json:"title"`
	Description     *string          `db:"description" json:"description"`
	Status          string           `db:"status" json:"status"`
	DueAt           *time.Time       `db:"due_at" json:"due_at"`
	LocationTrigger *json.RawMessage `db:"location_trigger" json:"location_trigger"`
	Source          *string          `db:"source" json:"source"`
	ExternalID      *string          `db:"external_id" json:"external_id"`
	CreatedAt       time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time        `db:"updated_at" json:"updated_at"`
}

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetAll() ([]Task, error) {
	tasks := []Task{}
	query := `SELECT * FROM tasks ORDER BY created_at DESC`
	err := s.db.Select(&tasks, query)
	return tasks, err
}

func (s *Store) GetByID(id string) (*Task, error) {
	var t Task
	query := `SELECT * FROM tasks WHERE id = $1`
	err := s.db.Get(&t, query, id)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *Store) Update(t *Task) error {
	query := `
		UPDATE tasks 
		SET title = :title, description = :description, status = :status, 
		    due_at = :due_at, location_trigger = :location_trigger, 
		    updated_at = NOW()
		WHERE id = :id`
	_, err := s.db.NamedExec(query, t)
	return err
}

func (s *Store) Delete(id string) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := s.db.Exec(query, id)
	return err
}

func (s *Store) Create(t *Task) error {
	query := `
		INSERT INTO tasks (title, description, status, due_at, location_trigger, source, external_id)
		VALUES (:title, :description, :status, :due_at, :location_trigger, :source, :external_id)
		RETURNING id, created_at, updated_at`
	
	rows, err := s.db.NamedQuery(query, t)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
	}
	return nil
}

func (s *Store) CreateRaw(data map[string]interface{}) error {
	query := `
		INSERT INTO tasks (title, description, status, due_at, source, external_id)
		VALUES (:title, :description, :status, :due_at, :source, :external_id)
		ON CONFLICT (external_id) DO NOTHING`
	_, err := s.db.NamedExec(query, data)
	return err
}
