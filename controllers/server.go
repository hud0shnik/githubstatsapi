package controllers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	api "github.com/hud0shnik/githubstatsapi/api"
	api2 "github.com/hud0shnik/githubstatsapi/api/v2"
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
	router.Get("/api/user", api.User)
	router.Get("/api/repo", api.Repo)
	router.Get("/api/commits", api.Commits)

	// Маршруты v2
	router.Get("/api/v2/user", api2.User)
	router.Get("/api/v2/repo", api2.Repo)
	router.Get("/api/v2/commits", api2.Commits)

	s.router = router

}
