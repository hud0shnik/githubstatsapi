package handler2

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
type userCommits struct {
	Success  bool   `json:"success"`
	Error    string `json:"error"`
	Date     string `json:"date"`
	Username string `json:"username"`
	Commits  int    `json:"commits"`
	Color    int    `json:"color"`
}

// Функция получения коммитов
func getCommits(username string, date string) userCommits {

	// Если поле даты пустое, функция поставит сегодняшнее число
	if date == "" {
		date = string(time.Now().Format("2006-01-02"))
	}

	// Формирование и исполнение запроса
	resp, err := http.Get("https://github.com/" + username + "?tab=overview&from=" + date)
	if err != nil {
		return userCommits{
			Error: "can't reach github.com",
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
	result := userCommits{
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
		result.Error = "not found"
	}

	return result
}

// Роут "/commits"
func Commits(w http.ResponseWriter, r *http.Request) {

	// Передача в заголовок респонса типа данных
	w.Header().Set("Content-Type", "application/json")

	// Получение параметра id из реквеста
	id := r.URL.Query().Get("id")

	// Проверка на наличие параметра
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(apiError{Error: "please insert user id"})
		w.Write(json)
		return
	}

	// Получение статистики и перевод в json
	result := getCommits(id, r.URL.Query().Get("date"))
	jsonResp, err := json.Marshal(result)

	// Обработчик ошибок
	switch {
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(apiError{Error: "internal server error"})
		w.Write(json)
		log.Printf("json.Marshal error: %s", err)
	case result.Error == "not found":
		w.WriteHeader(http.StatusNotFound)
		json, _ := json.Marshal(apiError{Error: "not found"})
		w.Write(json)
	case !result.Success:
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(apiError{Error: result.Error})
		w.Write(json)
	default:
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	}
}
