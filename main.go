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

// Структура ошибки
type ErrorResponse struct {
	Description string `json:"error"`
}

// Структура для хранения полной информации о пользователе
type UserInfo struct {
	Username      string `json:"username"`
	Name          string `json:"name"`
	Followers     string `json:"followers"`
	Following     string `json:"following"`
	Repositories  string `json:"repositories"`
	Packages      string `json:"packages"`
	Stars         string `json:"stars"`
	Contributions string `json:"contributions"`
	Status        string `json:"status"`
	Avatar        string `json:"avatar"`
}

// Структура для хранения информации о репозитории
type RepoInfo struct {
	Username string `json:"username"`
	Reponame string `json:"reponame"`
	Commits  string `json:"commits"`
	Branches string `json:"branches"`
	Tags     string `json:"tags"`
	Stars    string `json:"stars"`
	Watching string `json:"watching"`
	Forks    string `json:"forks"`
}

// Структура для хранения информации о коммитах
type UserCommits struct {
	Date     string `json:"date"`
	Username string `json:"username"`
	Commits  int    `json:"commits"`
	Color    int    `json:"color"`
}

// Функция поиска. Возвращает искомое значение и индекс последнего символа
func findWithIndex(str, subStr, stopChar string, start int) (string, int) {

	// Обрезка левой границы поиска
	str = str[start:]

	// Проверка на существование нужной строки
	if strings.Contains(str, subStr) {

		// Поиск индекса начала нужной строки
		left := strings.Index(str, subStr) + len(subStr)

		// Поиск правой границы
		right := left + strings.Index(str[left:], stopChar)

		// Обрезка и вывод результата
		return str[left:right], right + start
	}

	return "", 0
}

// Облегчённая функция поиска. Возвращает только искомое значение
func find(str, subStr, stopChar string) string {

	// Проверка на существование нужной строки
	if strings.Contains(str, subStr) {

		// Обрезка левой части
		str = str[strings.Index(str, subStr)+len(subStr):]

		// Обрезка правой части и вывод результата
		return str[:strings.Index(str, stopChar)]
	}

	return ""
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

	// Запись html в файл для тестирования
	/*if err := os.WriteFile("sample.html", []byte(pageStr), 0666); err != nil {
		log.Fatal(err)
	}*/

	// Структура, которую будет возвращать функция
	result := UserCommits{
		Date:     date,
		Username: username,
	}

	// Индекс ячейки с нужной датой
	i := strings.Index(pageStr, "data-date=\""+date)

	// Поиск и запись информации
	if i != -1 {
		pageStr = pageStr[i-22:]
		result.Commits, _ = strconv.Atoi(find(pageStr, "data-count=\"", "\""))
		result.Color, _ = strconv.Atoi(find(pageStr, "data-level=\"", "\""))
	}

	return result
}

// Функция получения информации о репозитории
func getRepoInfo(username string, reponame string) RepoInfo {

	// Формирование и исполнение запроса
	resp, err := http.Get("https://github.com/" + username + "/" + reponame)
	if err != nil {
		return RepoInfo{}
	}

	// Запись респонса
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// HTML полученной страницы в формате string
	pageStr := string(body)[30000:]

	// Запись html в файл для тестирования
	/*if err := os.WriteFile("sampleRepo.html", []byte(pageStr), 0666); err != nil {
		log.Fatal(err)
	}*/

	// Проверка на репозиторий
	if !strings.Contains(pageStr, "name=\"selected-link\" value=\"repo_source\"") {
		return RepoInfo{}
	}

	// Структура, которую будет возвращать функция
	result := RepoInfo{
		Username: username,
		Reponame: reponame,
	}

	// Индекс конца последней найденной строки
	left := 0

	// Ветки
	result.Branches, left = findWithIndex(pageStr, "01-1.5 0z\"></path>\n</svg>\n          <strong>", "<", left)

	// Теги
	result.Tags, left = findWithIndex(pageStr, "000-2z\"></path>\n</svg>\n        <strong>", "<", left)

	// Коммиты
	result.Commits, left = findWithIndex(pageStr, "class=\"d-none d-sm-inline\">\n                    <strong>", "<", left)

	// Звезды
	result.Stars, left = findWithIndex(pageStr, "94v.001z\"></path>\n</svg>\n    <strong>", "<", left)

	// Просмотры
	result.Watching, left = findWithIndex(pageStr, " 000 4z\"></path>\n</svg>\n    <strong>", "<", left)

	// Форки
	result.Forks, _ = findWithIndex(pageStr, "5.75.75 0 000 1.5z\"></path>\n</svg>\n    <strong>", "<", left)

	return result
}

// Функция получения информации о пользователе
func getUserInfo(username string) UserInfo {

	// Формирование и исполнение запроса
	resp, err := http.Get("https://github.com/" + username)
	if err != nil {
		return UserInfo{}
	}

	// Запись респонса
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// HTML полученной страницы в формате string
	pageStr := string(body)[55000:]

	// Запись html в файл для тестирования
	/*if err := os.WriteFile("sample.html", []byte(pageStr), 0666); err != nil {
		log.Fatal(err)
	}*/

	// Проверка на страницу пользователя и доступ к коммитам
	if !strings.Contains(pageStr, "p-nickname vcard-username d-block") || strings.Contains(pageStr, "class=\"octicon octicon-lock\">") {
		return UserInfo{}
	}

	// Структура, которую будет возвращать функция
	result := UserInfo{
		Username: username,
	}

	// Индекс конца последней найденной строки
	left := 0

	// Репозитории
	result.Repositories, left = findWithIndex(pageStr, "Repositories\n    <span title=\"", "\"", left)

	// Пакеты
	result.Packages, left = findWithIndex(pageStr, "Packages\n      <span title=\"", "\"", left)

	// Поставленные звезды
	result.Stars, left = findWithIndex(pageStr, "Stars\n    <span title=\"", "\"", left)

	// Ссылка на аватар
	result.Avatar, left = findWithIndex(pageStr, "r color-bg-default\" src=\"", "?", left)

	// Статус
	result.Status, left = findWithIndex(pageStr, "        <div>", "<", left)

	// Имя пользователя
	result.Name, left = findWithIndex(pageStr, "\"name\">\n          ", "\n", left)

	// Подписчики
	result.Followers, left = findWithIndex(pageStr, "lt\">", "<", left)

	// Подписки
	result.Following, left = findWithIndex(pageStr, "lt\">", "<", left)

	// Контрибуции за год
	result.Contributions, _ = findWithIndex(pageStr, "l mb-2\">\n      ", "\n", left)

	return result
}

// Функция отправки коммитов
func sendCommits(writer http.ResponseWriter, request *http.Request) {

	// Заголовок, определяющий тип данных респонса
	writer.Header().Set("Content-Type", "application/json")

	// Обработка данных и вывод результата
	json.NewEncoder(writer).Encode(getCommits(mux.Vars(request)["id"], mux.Vars(request)["date"]))
}

// Функция отправки информации о репозитории
func sendRepoInfo(writer http.ResponseWriter, request *http.Request) {

	// Заголовок, определяющий тип данных респонса
	writer.Header().Set("Content-Type", "application/json")

	// Обработка данных и вывод результата
	json.NewEncoder(writer).Encode(getRepoInfo(mux.Vars(request)["id"], mux.Vars(request)["repo"]))
}

// Функция отправки информации о пользователе
func sendUserInfo(writer http.ResponseWriter, request *http.Request) {

	// Заголовок, определяющий тип данных респонса
	writer.Header().Set("Content-Type", "application/json")

	// Обработка данных и вывод результата
	json.NewEncoder(writer).Encode(getUserInfo(mux.Vars(request)["id"]))
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

	router.HandleFunc("/user/{id}", sendUserInfo).Methods("GET")
	router.HandleFunc("/user/{id}/", sendUserInfo).Methods("GET")

	router.HandleFunc("/repo/{id}/{repo}", sendRepoInfo).Methods("GET")
	router.HandleFunc("/repo/{id}/{repo}/", sendRepoInfo).Methods("GET")

	// Запуск API

	// Для Heroku
	//log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))

	// Для локалхоста (127.0.0.1:8080/)
	log.Fatal(http.ListenAndServe(":8080", router))
}
