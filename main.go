package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Структура для храния полной информации о пользователе
type UserStats struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Stars    int    `json:"stars"`
}

// Структура для храния информации о коммитах
type UserCommits struct {
	Date     string `json:"date"`
	Username string `json:"username"`
	Commits  int    `json:"commits"`
	Color    int    `json:"color"`
}

// Функция получения статистики с сайта
func getStats(username string) UserStats {
	// Формирование и исполнение запроса
	resp, err := http.Get("https://github.com/" + username)
	if err != nil {
		return UserStats{}
	}

	// Запись респонса
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// HTML полученной страницы в формате string
	pageStr := string(body)

	// Структура, которую будет возвращать функция
	result := UserStats{
		Username: username,
	}

	// Проверка на страницу пользователя
	if !strings.Contains(pageStr, "p-nickname vcard-username d-block") {
		return result
	}

	// Обрезка ненужных частей страницы
	pageStr = pageStr[100000:195000]

	// Поиск информации о звездах
	left := strings.Index(pageStr, "Stars\n    <span title=") + 23

	// Если звезды есть, считывает их количество и записывает
	if left > 23 {
		right := left
		for ; pageStr[right] != '"'; right++ {
			// Доводит pageStr[right] до символа '"'
		}

		// Запись звезд в результат
		result.Stars, _ = strconv.Atoi(pageStr[left:right])
	}

	// Поиск имени пользователя
	left = strings.Index(pageStr, "itemprop=\"n") + 16

	// Если имя найдено, считывает его и записывает
	if left > 16 {
		right := left
		for ; pageStr[right] != '<'; right++ {
			// Доводит pageStr[right] до символа '<'
		}

		// Запись имени
		result.Name = pageStr[left+11 : right-9]
	}

	// Поиск ссылки на аватар
	left = strings.Index(pageStr, "https://avatars.githubusercontent.com/u")

	// Если ссылка найдена, считывает её и записывает
	if left != -1 {
		right := left
		for ; pageStr[right] != '?'; right++ {
			// Доводит pageStr[right] до символа '?'
		}

		// Запись ссылки
		result.Avatar = pageStr[left:right]
	}

	return result
}

// Функция получения коммитов
func getCommits(username string, date string) UserCommits {
	// Формирование и исполнение запроса
	resp, err := http.Get("https://github.com/" + username)
	if err != nil {
		return UserCommits{}
	}

	// Запись респонса
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// HTML полученной страницы в формате string
	pageStr := string(body)

	// Если поле даты пустое, функция поставит сегодняшнее число
	if date == "" {
		date = string(time.Now().Format("2006-01-02"))
	}

	// Структура, которую будет возвращать функция
	result := UserCommits{
		Date:     date,
		Username: username,
	}

	// Проверка на страницу пользователя
	if !strings.Contains(pageStr, "p-nickname vcard-username d-block") {
		return result
	}

	// Обрезка ненужных частей страницы
	pageStr = pageStr[100000:]

	// Указатель на ячейку нужной даты
	i := strings.Index(pageStr, "data-date=\""+date)
	// Проверка на существование нужной ячейки
	if i != -1 {
		for ; pageStr[i] != '<'; i-- {
			// Доводит i до начала тега ячейки
		}

		// Получение параметров ячейки
		values := strings.FieldsFunc(pageStr[i:i+150], func(r rune) bool {
			return r == '"'
		})

		// Запись и обработка нужной информации
		result.Color, _ = strconv.Atoi(values[19])
		result.Commits, _ = strconv.Atoi(values[15])

	}

	return result
}

// Функция отправки коммитов
func sendCommits(writer http.ResponseWriter, request *http.Request) {
	// Заголовок, определяющий тип данных респонса
	writer.Header().Set("Content-Type", "application/json")

	// Обработка данных и вывод результата
	json.NewEncoder(writer).Encode(getCommits(mux.Vars(request)["id"], mux.Vars(request)["date"]))
}

// Функция отправки статистики
func sendStats(writer http.ResponseWriter, request *http.Request) {
	// Заголовок, определяющий тип данных респонса
	writer.Header().Set("Content-Type", "application/json")

	// Обработка данных и вывод результата
	json.NewEncoder(writer).Encode(getStats(mux.Vars(request)["id"]))
}

func main() {
	// Вывод времени начала работы
	fmt.Println("API Start: " + string(time.Now().Format("2006-01-02 15:04:05")))

	// Роутер
	router := mux.NewRouter()

	// Маршруты
	router.HandleFunc("/commits/{id}", sendCommits).Methods("GET")
	router.HandleFunc("/commits/{id}/", sendCommits).Methods("GET")
	router.HandleFunc("/commits/{id}/{date}", sendCommits).Methods("GET")
	router.HandleFunc("/commits/{id}/{date}/", sendCommits).Methods("GET")

	router.HandleFunc("/stats/{id}", sendStats).Methods("GET")
	router.HandleFunc("/stats/{id}/", sendStats).Methods("GET")

	// Запуск API

	// Для Heroku
	//log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))

	// Для локалхоста (127.0.0.1:8080/)
	log.Fatal(http.ListenAndServe(":8080", router))
}
