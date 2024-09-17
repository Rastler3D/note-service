package api

import (
	"encoding/json"
	"net/http"
	"todo-list/internal/logger"
	"todo-list/internal/models"
	"todo-list/internal/repository"
	"todo-list/internal/services"
)

type Handler struct {
	repo       repository.NoteRepository
	spellcheck services.SpellcheckService
}

func NewHandler(repo repository.NoteRepository, spellcheck services.SpellcheckService) *Handler {
	return &Handler{repo: repo, spellcheck: spellcheck}
}

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var note models.Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		logger.Error("Error decoding request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int64)
	note.UserID = userID

	logger.Info("Creating note for user %d", userID)

	correctedTitle, err := h.spellcheck.CorrectText(note.Title)
	if err != nil {
		logger.Error("Error correcting title: %v", err)
		http.Error(w, "Error correcting title: "+err.Error(), http.StatusInternalServerError)
		return
	}
	note.Title = correctedTitle

	correctedContent, err := h.spellcheck.CorrectText(note.Content)
	if err != nil {
		logger.Error("Error correcting content: %v", err)
		http.Error(w, "Error correcting content: "+err.Error(), http.StatusInternalServerError)
		return
	}
	note.Content = correctedContent

	err = h.repo.Create(&note)
	if err != nil {
		logger.Error("Error creating note: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Note created successfully: %d", note.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

func (h *Handler) GetNotes(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	logger.Info("Fetching notes for user %d", userID)

	notes, err := h.repo.GetByUserID(userID)
	if err != nil {
		logger.Error("Error fetching notes: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Info("Successfully fetched %d notes for user %d", len(notes), userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}
