package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"todo-list/internal/config"
	"todo-list/internal/logger"
	"todo-list/internal/repository"
	"todo-list/internal/server"
	"todo-list/internal/services"
)

func main() {
	logger.Init()

	cfg := config.NewConfig()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Error("Error opening database: %v", err)
		return
	}
	defer db.Close()
	srv := server.NewServer(cfg.ServerPort)
	repo, err := repository.NewNoteRepository(db)
	if err != nil {
		logger.Error("Error initializing database: %v", err)
		return
	}
	spellcheckService := services.NewSpellcheckService(cfg.SpellcheckerURL)
	srv.SetupRoutes(repo, spellcheckService)

	if err := srv.Start(); err != nil {
		logger.Error("Server error: %v", err)
	}
}
