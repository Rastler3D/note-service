package repository

import (
	"database/sql"
	"time"
	"todo-list/internal/logger"
	"todo-list/internal/models"
)

type NoteRepository interface {
	Create(note *models.Note) error
	GetByUserID(userID int64) ([]models.Note, error)
}

type DbNoteRepository struct {
	db *sql.DB
}

func NewNoteRepository(db *sql.DB) (NoteRepository, error) {
	err := initDatabase(db)
	if err != nil {
		return nil, err
	}
	return &DbNoteRepository{db: db}, nil
}

func initDatabase(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS notes (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_notes_user_id ON notes(user_id);
	`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	logger.Info("Notes table created or already exists")
	return nil
}

func (r *DbNoteRepository) Create(note *models.Note) error {
	query := `INSERT INTO notes (user_id, title, content, created_at) 
              VALUES ($1, $2, $3, $4) RETURNING id`

	now := time.Now()
	err := r.db.QueryRow(query, note.UserID, note.Title, note.Content, now).Scan(&note.ID)
	if err != nil {
		return err
	}

	note.CreatedAt = now

	return nil
}

func (r *DbNoteRepository) GetByUserID(userID int64) ([]models.Note, error) {
	query := `SELECT id, user_id, title, content, created_at
              FROM notes WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := []models.Note{}
	for rows.Next() {
		var note models.Note
		err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}
