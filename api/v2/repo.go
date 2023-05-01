package api2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
func getRepoInfoString(username, reponame string) (repoInfoString, error) {

	// Формирование и исполнение запроса
	resp, err := http.Get("https://github.com/" + username + "/" + reponame)
	if err != nil {
		return repoInfoString{}, fmt.Errorf("in http.Get: %w", err)

	}

	// Запись респонса
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// HTML полученной страницы в формате string
	pageStr := string(body)[20000:]

	// Запись html в файл для тестирования
	/*if err := os.WriteFile("sampleRepo.html", []byte(pageStr), 0666); err != nil {
		log.Fatal(err)
	}*/

	// Проверка на репозиторий
	if !strings.Contains(pageStr, "name=\"selected-link\" value=\"repo_source\"") {
		return repoInfoString{}, fmt.Errorf("not found")
	}

	// Структура, которую будет возвращать функция
	result := repoInfoString{
		Username: username,
		Reponame: reponame,
	}

	// Индекс конца последней найденной строки
	left := 0

	// Звезды
	result.Stars, left = findWithIndex(pageStr, "2.694Z\"></path>\n</svg>\n          <span class=\"text-bold\">", "<", left)

	// Форки
	result.Forks, left = findWithIndex(pageStr, "0Z\"></path>\n</svg>\n          <span class=\"text-bold\">", "<", left)

	// Ветки
	result.Branches, left = findWithIndex(pageStr, "0-1.5Z\"></path>\n</svg>\n          <strong>", "<", left)

	// Теги
	result.Tags, left = findWithIndex(pageStr, "0-2Z\"></path>\n</svg>\n        <strong>", "<", left)

	// Коммиты
	result.Commits, left = findWithIndex(pageStr, "class=\"d-none d-sm-inline\">\n                    <strong>", "<", left)

	// Просмотры
	result.Watching, _ = findWithIndex(pageStr, "10Z\"></path>\n</svg>\n    <strong>", "<", left)

	return result, nil

}

// Функция получения информации о репозитории
func getRepoInfo(username, reponame string) (repoInfo, error) {

	// Получение текстовой версии статистики
	resultStr, err := getRepoInfoString(username, reponame)
	if err != nil {
		return repoInfo{}, err
	}

	return repoInfo{
		Username: resultStr.Username,
		Reponame: resultStr.Reponame,
		Commits:  toInt(resultStr.Commits),
		Branches: toInt(resultStr.Branches),
		Tags:     toInt(resultStr.Tags),
		Stars:    toInt(resultStr.Stars),
		Watching: toInt(resultStr.Watching),
		Forks:    toInt(resultStr.Forks),
	}, nil

}

// Роут "/repo"
func Repo(w http.ResponseWriter, r *http.Request) {

	// Передача в заголовок респонса типа данных
	w.Header().Set("Content-Type", "application/json")

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
		result, err := getRepoInfoString(username, reponame)
		if err != nil {
			if err.Error() == "not found" {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
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
			log.Printf("json.Marshal error: %s", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)

	} else {

		// Получение статистики и перевод в json
		result, err := getRepoInfo(username, reponame)
		if err != nil {
			if err.Error() == "not found" {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
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
			log.Printf("json.Marshal error: %s", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)

	}

}
