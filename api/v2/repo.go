package api2

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hud0shnik/githubstatsapi/utils"
	"github.com/sirupsen/logrus"
)

// Структура для хранения информации о репозитории
type repoInfo struct {
	Username string `json:"username"`
	Reponame string `json:"reponame"`
	Commits  int    `json:"commits"`
	Branches int    `json:"branches"`
	Tags     int    `json:"tags"`
	Stars    int    `json:"stars"`
	Watching int    `json:"watching"`
	Forks    int    `json:"forks"`
}

// Структура для парсинга информации о репозитории
type repoInfoString struct {
	Username string `json:"username"`
	Reponame string `json:"reponame"`
	Commits  string `json:"commits"`
	Branches string `json:"branches"`
	Tags     string `json:"tags"`
	Stars    string `json:"stars"`
	Watching string `json:"watching"`
	Forks    string `json:"forks"`
}

// Функция получения информации о репозитории в формате строк
func getRepoInfoString(username, reponame string) (repoInfoString, int, error) {

	// Формирование и исполнение запроса
	resp, err := http.Get("https://github.com/" + username + "/" + reponame)
	if err != nil {
		return repoInfoString{}, http.StatusInternalServerError, fmt.Errorf("in http.Get: %w", err)
	}
	defer resp.Body.Close()

	// Проверка статускода
	if resp.StatusCode != 200 {
		return repoInfoString{}, resp.StatusCode,
			fmt.Errorf(resp.Status)
	}

	// Запись респонса
	body, _ := io.ReadAll(resp.Body)

	// HTML полученной страницы в формате string
	pageStr := string(body)[20000:]

	// Запись html в файл для тестирования
	/*if err := os.WriteFile("sampleRepo.html", []byte(pageStr), 0666); err != nil {
		logrus.Fatal(err)
	}*/

	// Проверка на репозиторий
	if !strings.Contains(pageStr, "name=\"selected-link\" value=\"repo_source\"") {
		return repoInfoString{}, http.StatusNotFound, fmt.Errorf("not found")
	}

	// Структура, которую будет возвращать функция
	result := repoInfoString{
		Username: username,
		Reponame: reponame,
	}

	// Индекс конца последней найденной строки
	left := 0

	// Звезды
	result.Stars, left = utils.FindWithIndex(pageStr, "2.694Z\"></path>\n</svg>\n          <span class=\"text-bold\">", "<", left)

	// Форки
	result.Forks, left = utils.FindWithIndex(pageStr, "0Z\"></path>\n</svg>\n          <span class=\"text-bold\">", "<", left)

	// Ветки
	result.Branches, left = utils.FindWithIndex(pageStr, "0-1.5Z\"></path>\n</svg>\n          <strong>", "<", left)

	// Теги
	result.Tags, left = utils.FindWithIndex(pageStr, "0-2Z\"></path>\n</svg>\n        <strong>", "<", left)

	// Коммиты
	result.Commits, left = utils.FindWithIndex(pageStr, "class=\"d-none d-sm-inline\">\n                    <strong>", "<", left)

	// Просмотры
	result.Watching, _ = utils.FindWithIndex(pageStr, "10Z\"></path>\n</svg>\n    <strong>", "<", left)

	return result, http.StatusOK, nil

}

// Функция получения информации о репозитории
func getRepoInfo(username, reponame string) (repoInfo, int, error) {

	// Получение текстовой версии статистики
	resultStr, statusCode, err := getRepoInfoString(username, reponame)
	if err != nil {
		return repoInfo{}, statusCode, err
	}

	return repoInfo{
		Username: resultStr.Username,
		Reponame: resultStr.Reponame,
		Commits:  utils.ToInt(resultStr.Commits),
		Branches: utils.ToInt(resultStr.Branches),
		Tags:     utils.ToInt(resultStr.Tags),
		Stars:    utils.ToInt(resultStr.Stars),
		Watching: utils.ToInt(resultStr.Watching),
		Forks:    utils.ToInt(resultStr.Forks),
	}, http.StatusOK, nil

}

// Роут "/repo"
func Repo(w http.ResponseWriter, r *http.Request) {

	// Установка заголовков
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// Получение имени пользователя и названия репозитория из реквеста
	username := r.URL.Query().Get("username")
	reponame := r.URL.Query().Get("reponame")

	// Проверка на наличие параметров
	if username == "" || reponame == "" {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(apiError{Error: "please insert username and reponame"})
		w.Write(json)
		return
	}

	// Проверка на тип
	if r.URL.Query().Get("type") == "string" {

		// Получение статистики
		result, statusCode, err := getRepoInfoString(username, reponame)
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

	} else {

		// Получение статистики и перевод в json
		result, statusCode, err := getRepoInfo(username, reponame)
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

}
