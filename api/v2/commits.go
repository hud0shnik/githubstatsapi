package api2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Структура для хранения информации о коммитах
type userCommits struct {
	Date     string `json:"date"`
	Username string `json:"username"`
	Commits  int    `json:"commits"`
	Color    int    `json:"color"`
}

// Функция получения коммитов
func getCommits(username string, date string) (userCommits, int, error) {

	// Если поле даты пустое, функция поставит сегодняшнее число
	if date == "" {
		date = string(time.Now().Format("2006-01-02"))
	}

	// Формирование и исполнение запроса
	resp, err := http.Get("https://github.com/" + username + "?tab=overview&from=" + date)
	if err != nil {
		return userCommits{}, http.StatusInternalServerError,
			fmt.Errorf("in http.Get: %w", err)
	}
	defer resp.Body.Close()

	// Проверка статускода
	if resp.StatusCode != 200 {
		return userCommits{}, resp.StatusCode,
			fmt.Errorf(resp.Status)
	}

	// Запись респонса
	body, _ := ioutil.ReadAll(resp.Body)

	// HTML полученной страницы в формате string
	pageStr := string(body)[100000:]

	// Запись html в файл для тестирования
	/*if err := os.WriteFile("sample.html", []byte(pageStr), 0666); err != nil {
		logrus..Fatal(err)
	}*/

	// Структура, которую будет возвращать функция
	result := userCommits{
		Date:     date,
		Username: username,
	}

	// Индекс ячейки с нужной датой
	i := strings.Index(pageStr, "data-date=\""+date)

	// Проверка на наличие ячейки
	if i == -1 {
		return userCommits{}, http.StatusNotFound,
			fmt.Errorf("not found")
	}

	// Запись данных
	pageStr = pageStr[i-22:]
	result.Color, _ = strconv.Atoi(find(pageStr, "data-level=\"", "\""))
	result.Commits, _ = strconv.Atoi(find(pageStr, "\">", " "))

	return result, http.StatusOK, nil

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

	// Получение статистики
	result, statusCode, err := getCommits(id, r.URL.Query().Get("date"))
	if err != nil {
		w.WriteHeader(statusCode)
		json, _ := json.Marshal(apiError{Error: err.Error()})
		w.Write(json)
		return
	}

	// Перевод в json
	jsonResp, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(apiError{Error: "internal server error"})
		w.Write(json)
		logrus.Printf("json.Marshal error: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)

}
