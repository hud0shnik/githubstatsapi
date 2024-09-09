package controller

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hud0shnik/githubstatsapi/internal/handler"
	handlerv2 "github.com/hud0shnik/githubstatsapi/internal/handler/v2"
)

// Server - структура сервера
type Server struct {
	basePath       string
	requestTimeout time.Duration
	router         http.Handler
	Server         *http.Server
}

// NewServer создаёт новый сервер
func NewServer(config *Config) *Server {

	s := &Server{
		basePath:       config.BasePath,
		requestTimeout: config.RequestTimeout,
	}

	s.NewRouter()

	s.Server = &http.Server{
		Addr:              config.ServerPort,
		Handler:           s.router,
		ReadTimeout:       config.RequestTimeout,
		ReadHeaderTimeout: config.RequestTimeout,
	}

	return s
}

// NewRouter создаёт новый роутер
func (s *Server) NewRouter() {

	// Роутер
	router := chi.NewRouter()

	// Маршруты
	router.Get("/api/user", handler.User)
	router.Get("/api/repo", handler.Repo)
	router.Get("/api/commits", handler.Commits)

	// Маршруты v2
	router.Get("/api/v2/user", handlerv2.User)
	router.Get("/api/v2/repo", handlerv2.Repo)
	router.Get("/api/v2/commits", handlerv2.Commits)

	s.router = router

}
