package main

import (
	"fmt"
	api "gitAPI/api"
	api2 "gitAPI/api/v2"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	// Вывод времени начала работы
	fmt.Println("API Start: " + string(time.Now().Format("2006-01-02 15:04:05")))
	fmt.Println("Port:\t" + os.Getenv("PORT"))

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
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))

}
