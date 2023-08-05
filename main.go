package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/hud0shnik/githubstatsapi/api"
	api2 "github.com/hud0shnik/githubstatsapi/api/v2"
	"github.com/sirupsen/logrus"
)

func main() {

	// Настройка логгера
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.DateTime,
	})

	// Вывод времени начала работы
	logrus.Info("API Start")
	logrus.Info("Port: " + os.Getenv("PORT"))

	// Роутер
	router := chi.NewRouter()

	// Маршруты

	router.Get("/api/user", api.User)
	router.Get("/api/repo", api.Repo)
	router.Get("/api/commits", api.Commits)

	router.Get("/api/v2/user", api2.User)
	router.Get("/api/v2/repo", api2.Repo)
	router.Get("/api/v2/commits", api2.Commits)

	// Запуск API
	logrus.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))

}
