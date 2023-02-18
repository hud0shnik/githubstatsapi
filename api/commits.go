package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Структура для хранения информации о коммитах
type UserCommits struct {
	Success  bool   `json:"success"`
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
			Error: "http get error",
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
		result.Success = true
		pageStr = pageStr[i-22:]
		result.Color, _ = strconv.Atoi(find(pageStr, "data-level=\"", "\""))
		result.Commits, _ = strconv.Atoi(find(pageStr, "\">", " "))
	} else {
		result.Error = "commits not found"
	}

	return result
}

// Роут "/commits"
func Commits(w http.ResponseWriter, r *http.Request) {

	// Получение параметра id из реквеста
	id := r.URL.Query().Get("id")

	// Если параметра нет, отправка ошибки
	if id == "" {
		http.NotFound(w, r)
		return
	}

	// Формирование заголовка респонса по статускоду
	w.WriteHeader(http.StatusOK)

	// Передача в заголовок респонса типа данных
	w.Header().Set("Content-Type", "application/json")

	// Получение параметра даты из реквеста
	date := r.URL.Query().Get("date")

	// Получение и запись статистики
	resp := GetCommits(id, date)

	// Форматирование структуры в json и отправка пользователю
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("json.Marshal error: %s", err)
	} else {
		w.Write(jsonResp)
	}
}
