package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"todo-list/internal/logger"
	"todo-list/internal/models"
	"todo-list/internal/server"
)

type mockDB struct{}

func (m *mockDB) Create(note *models.Note) error {
	note.ID = 1
	return nil
}

func (m *mockDB) GetByUserID(userID int64) ([]models.Note, error) {
	return []models.Note{{ID: 1, UserID: userID, Title: "Test Note", Content: "This is a test note"}}, nil
}

type mockSpellchecker struct{}

func (m *mockSpellchecker) CorrectText(text string) (string, error) {
	return text, nil
}

type mockLogger struct {
	infoLogs  []string
	errorLogs []string
}

func (m *mockLogger) Info(format string, v ...interface{}) {
	m.infoLogs = append(m.infoLogs, fmt.Sprintf(format, v...))
}

func (m *mockLogger) Error(format string, v ...interface{}) {
	m.errorLogs = append(m.errorLogs, fmt.Sprintf(format, v...))
}

func TestCreateNote(t *testing.T) {
	serverPort := "8080"
	repo := &mockDB{}
	spellcheck := &mockSpellchecker{}
	mockLog := &mockLogger{}
	logger.InfoLogger = mockLog
	logger.ErrorLogger = mockLog

	srv := server.NewServer(serverPort)
	srv.SetupRoutes(repo, spellcheck)

	note := models.Note{Title: "Test Note", Content: "This is a test note"}
	body, _ := json.Marshal(note)

	req := httptest.NewRequest("POST", "/notes", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "token1")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	srv.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var createdNote models.Note
	json.Unmarshal(w.Body.Bytes(), &createdNote)

	if createdNote.ID != 1 {
		t.Errorf("Expected note ID 1, got %d", createdNote.ID)
	}

	// Check logs
	expectedInfoLog := "Creating note for user 1"
	found := false
	for _, log := range mockLog.infoLogs {
		if strings.Contains(log, expectedInfoLog) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected log message not found: %s", expectedInfoLog)
	}
}

func TestGetNotes(t *testing.T) {
	serverPort := "8080"
	repo := &mockDB{}
	spellcheck := &mockSpellchecker{}
	mockLog := &mockLogger{}
	logger.InfoLogger = mockLog
	logger.ErrorLogger = mockLog

	srv := server.NewServer(serverPort)
	srv.SetupRoutes(repo, spellcheck)

	req := httptest.NewRequest("GET", "/notes", nil)
	req.Header.Set("Authorization", "token1")

	w := httptest.NewRecorder()
	srv.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var notes []models.Note
	json.Unmarshal(w.Body.Bytes(), &notes)

	if len(notes) != 1 {
		t.Errorf("Expected 1 note, got %d", len(notes))
	}

	if notes[0].Title != "Test Note" {
		t.Errorf("Expected note title 'Test Note', got '%s'", notes[0].Title)
	}

	// Check logs
	expectedInfoLog := "Fetching notes for user 1"
	found := false
	for _, log := range mockLog.infoLogs {
		if strings.Contains(log, expectedInfoLog) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected log message not found: %s", expectedInfoLog)
	}
}

func TestCreateNoteInvalidJSON(t *testing.T) {
	serverPort := "8080"
	repo := &mockDB{}
	spellcheck := &mockSpellchecker{}
	mockLog := &mockLogger{}
	logger.InfoLogger = mockLog
	logger.ErrorLogger = mockLog

	srv := server.NewServer(serverPort)
	srv.SetupRoutes(repo, spellcheck)

	invalidJSON := []byte(`{"title": "Test Note", "content":}`) // Invalid JSON

	req := httptest.NewRequest("POST", "/notes", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Authorization", "token1")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	srv.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Check logs
	expectedErrorLog := "Error decoding request body"
	found := false
	for _, log := range mockLog.errorLogs {
		if strings.Contains(log, expectedErrorLog) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected log message not found: %s", expectedErrorLog)
	}
}
