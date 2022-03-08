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
type UserInfo struct {
	Username      string `json:"username"`
	Name          string `json:"name"`
	Followers     string `json:"followers"`
	Following     string `json:"following"`
	Repositories  string `json:"repositories"`
	Packages      string `json:"packages"`
	Stars         string `json:"stars"`
	Contributions string `json:"contributions"`
	Avatar        string `json:"avatar"`
}

// Структура для храния информации о коммитах
type UserCommits struct {
	Date     string `json:"date"`
	Username string `json:"username"`
	Commits  int    `json:"commits"`
	Color    int    `json:"color"`
}

// Функция поиска, возвращает искомое значение и индекс конца
func find(str string, subStr string, char byte) (string, int) {
	// Поиск индекса начала нужной строки
	left := strings.Index(str, subStr) + len(subStr)

	// Проверка на существование нужной строки
	if left > len(subStr)-1 {

		// Крайняя часть искомой строки
		right := left

		for ; str[right] != char; right++ {
			// Доводит str[right] до символа char
		}
		return str[left:right], right
	}

	return "", 0
}

// Функция получения информации о пользователе
func getInfo(username string) UserInfo {

	// Формирование и исполнение запроса
	resp, err := http.Get("https://github.com/" + username)
	if err != nil {
		return UserInfo{}
	}

	// Запись респонса
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// HTML полученной страницы в формате string
	pageStr := string(body)

	// Структура, которую будет возвращать функция
	result := UserInfo{
		Username: username,
	}

	// Проверка на страницу пользователя
	if !strings.Contains(pageStr, "p-nickname vcard-username d-block") {
		return result
	}

	// Уберает лишнюю часть
	pageStr = pageStr[:strings.Index(pageStr, "js-calendar-graph-svg")]

	// Индекс конца последней найденной строки
	i := 0

	/* -----------------------------------------------------------
	# Далее поиск нужной информации и запись в результат		 #
	# после каждого поиска тело сайта обрезается для оптимизации #
	------------------------------------------------------------ */

	// Репозитории
	result.Repositories, i = find(pageStr, "Repositories\n    <span title=\"", '"')
	pageStr = pageStr[i:]

	// Пакеты
	result.Packages, i = find(pageStr, "Packages\n      <span title=\"", '"')
	pageStr = pageStr[i:]

	// Поставленные звезды
	result.Stars, i = find(pageStr, "Stars\n    <span title=\"", '"')
	pageStr = pageStr[i:]

	// Ссылка на аватар
	result.Avatar, i = find(pageStr, "r color-bg-default\" src=\"", '?')
	pageStr = pageStr[i:]

	// Имя пользователя
	result.Name, i = find(pageStr, "\"name\">\n          ", '\n')
	pageStr = pageStr[i:]

	// Подписчики
	result.Followers, i = find(pageStr, "lt\">", '<')
	pageStr = pageStr[i:]

	// Подписки
	result.Following, i = find(pageStr, "lt\">", '<')
	pageStr = pageStr[i:]

	// Контрибуции за год
	result.Contributions, _ = find(pageStr, "l mb-2\">\n      ", '\n')

	return result
}

// Функция получения коммитов
func getCommits(username string, date string) UserCommits {

	// Если поле даты пустое, функция поставит сегодняшнее число
	if date == "" {
		date = string(time.Now().Format("2006-01-02"))
	}

	// Формирование и исполнение запроса
	resp, err := http.Get("https://github.com/" + username + "?tab=overview&from=" + date)
	if err != nil {
		return UserCommits{}
	}

	// Запись респонса
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// HTML полученной страницы в формате string
	pageStr := string(body)[100000:]

	// Структура, которую будет возвращать функция
	result := UserCommits{
		Date:     date,
		Username: username,
	}

	// Индекс ячейки с нужной датой
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

		// Запись нужной информации
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

// Функция отправки информации
func sendInfo(writer http.ResponseWriter, request *http.Request) {

	// Заголовок, определяющий тип данных респонса
	writer.Header().Set("Content-Type", "application/json")

	// Обработка данных и вывод результата
	json.NewEncoder(writer).Encode(getInfo(mux.Vars(request)["id"]))
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

	router.HandleFunc("/info/{id}", sendInfo).Methods("GET")
	router.HandleFunc("/info/{id}/", sendInfo).Methods("GET")

	// Запуск API

	// Для Heroku
	//log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))

	// Для локалхоста (127.0.0.1:8080/)
	log.Fatal(http.ListenAndServe(":8080", router))
}
