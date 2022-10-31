package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Структура для хранения информации о коммитах
type UserCommits struct {
	Error    string `json:"error"`
	Date     string `json:"date"`
	Username string `json:"username"`
	Commits  int    `json:"commits"`
	Color    int    `json:"color"`
}

// Функция получения коммитов
func GetCommits(username string, date string) UserCommits {

	// Если поле даты пустое, функция поставит сегодняшнее число
	if date == "" {
		date = string(time.Now().Format("2006-01-02"))
	}

	// Формирование и исполнение запроса
	resp, err := http.Get("https://github.com/" + username + "?tab=overview&from=" + date)
	if err != nil {
		return UserCommits{
			Error: "http.Get error",
		}
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
	} else {
		result.Error = "commits not found"
	}

	return result
}

// Роут "/commits"
func Commits(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.NotFound(w, r)
		return
	}
	date := r.URL.Query().Get("date")
	resp := GetCommits(id, date)

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Print("Error: ", err)
	} else {
		w.Write(jsonResp)
	}
}
