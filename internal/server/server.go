package server

import (
	"net/http"
	"todo-list/internal/api"
	"todo-list/internal/logger"
	"todo-list/internal/middleware"
	"todo-list/internal/repository"
	"todo-list/internal/services"
)

type Server struct {
	serverPort string
	Mux        *http.ServeMux
}

func NewServer(serverPort string) *Server {
	return &Server{
		serverPort: serverPort,
		Mux:        http.NewServeMux(),
	}
}

func (s *Server) SetupRoutes(repo repository.NoteRepository, spellcheckService services.SpellcheckService) {
	handler := api.NewHandler(repo, spellcheckService)

	s.Mux.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Received request: %s %s", r.Method, r.URL.Path)
		switch r.Method {
		case http.MethodPost:
			middleware.AuthMiddleware(handler.CreateNote).ServeHTTP(w, r)
		case http.MethodGet:
			middleware.AuthMiddleware(handler.GetNotes).ServeHTTP(w, r)
		default:
			logger.Error("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (s *Server) Start() error {
	logger.Info("Starting server on port %s", s.serverPort)
	return http.ListenAndServe(":"+s.serverPort, s.Mux)
}
