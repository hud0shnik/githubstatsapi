package main

import (
	"net/http"
	"os"
	"time"

	"github.com/hud0shnik/githubstatsapi/api"
	api2 "github.com/hud0shnik/githubstatsapi/api/v2"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
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
	router := mux.NewRouter()

	// Маршруты

	router.HandleFunc("/api/commits", api.Commits).Methods("GET")
	router.HandleFunc("/api/v2/commits", api2.Commits).Methods("GET")

	router.HandleFunc("/api/user", api.User).Methods("GET")
	router.HandleFunc("/api/v2/user", api2.User).Methods("GET")

	router.HandleFunc("/api/repo", api.Repo).Methods("GET")
	router.HandleFunc("/api/v2/repo", api2.Repo).Methods("GET")

	// Запуск API
	logrus.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))

}
