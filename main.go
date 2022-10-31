package main

import (
	"encoding/json"
	"fmt"
	api "gitAPI/api"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// Функция отправки коммитов
func sendCommits(writer http.ResponseWriter, request *http.Request) {

	// Заголовок, определяющий тип данных респонса
	writer.Header().Set("Content-Type", "application/json")

	// Обработка данных и вывод результата
	json.NewEncoder(writer).Encode(api.GetCommits(request.URL.Query().Get("id"), request.URL.Query().Get("date")))
}

// Функция отправки информации о репозитории
func sendRepoInfo(writer http.ResponseWriter, request *http.Request) {

	// Заголовок, определяющий тип данных респонса
	writer.Header().Set("Content-Type", "application/json")

	// Обработка данных и вывод результата
	json.NewEncoder(writer).Encode(api.GetRepoInfo(request.URL.Query().Get("username"), request.URL.Query().Get("reponame")))
}

// Функция отправки информации о пользователе
func sendUserInfo(writer http.ResponseWriter, request *http.Request) {

	// Заголовок, определяющий тип данных респонса
	writer.Header().Set("Content-Type", "application/json")

	// Обработка данных и вывод результата
	json.NewEncoder(writer).Encode(api.GetUserInfo(request.URL.Query().Get("id")))
}

func main() {

	// Вывод времени начала работы
	fmt.Println("API Start: " + string(time.Now().Format("2006-01-02 15:04:05")))
	fmt.Println("Port:\t" + os.Getenv("PORT"))

	// Роутер
	router := mux.NewRouter()

	// Маршруты

	router.HandleFunc("/api/commits", sendCommits).Methods("GET")

	router.HandleFunc("/api/user", sendUserInfo).Methods("GET")

	router.HandleFunc("/api/repo", sendRepoInfo).Methods("GET")

	// Запуск API
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))

}
