package handler2

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Структура для хранения информации о репозитории
type RepoInfo struct {
	Success  bool   `json:"success"`
	Error    string `json:"error"`
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
type RepoInfoString struct {
	Success  bool   `json:"success"`
	Error    string `json:"error"`
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
func GetRepoInfoString(username, reponame string) RepoInfoString {

	// Формирование и исполнение запроса
	resp, err := http.Get("https://github.com/" + username + "/" + reponame)
	if err != nil {
		return RepoInfoString{
			Error: "can't reach github.com",
		}
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
		return RepoInfoString{
			Error: "not found",
		}
	}

	// Структура, которую будет возвращать функция
	result := RepoInfoString{
		Success:  true,
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

// Функция получения информации о репозитории
func GetRepoInfo(username, reponame string) RepoInfo {

	// Получение текстовой версии статистики
	resultStr := GetRepoInfoString(username, reponame)

	// Проверка на ошибки при парсинге
	if !resultStr.Success {
		return RepoInfo{
			Success: false,
			Error:   resultStr.Error,
		}
	}

	return RepoInfo{
		Success:  resultStr.Success,
		Error:    resultStr.Error,
		Username: resultStr.Username,
		Reponame: resultStr.Reponame,
		Commits:  toInt(resultStr.Commits),
		Branches: toInt(resultStr.Branches),
		Tags:     toInt(resultStr.Tags),
		Stars:    toInt(resultStr.Stars),
		Watching: toInt(resultStr.Watching),
		Forks:    toInt(resultStr.Forks),
	}
}

// Роут "/repo"
func Repo(w http.ResponseWriter, r *http.Request) {

	// Передача в заголовок респонса типа данных
	w.Header().Set("Content-Type", "application/json")

	// Получение параметра имени пользователя и названия репозитория из реквеста
	username := r.URL.Query().Get("username")
	reponame := r.URL.Query().Get("reponame")

	// Если параметра нет, отправка ошибки
	if username == "" || reponame == "" {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(ApiError{Error: "please insert username and reponame"})
		w.Write(json)
		return
	}

	// Проверка на тип, получение статистики, форматирование и отправка
	if r.URL.Query().Get("type") == "string" {
		result := GetRepoInfoString(username, reponame)
		jsonResp, err := json.Marshal(result)
		switch {
		case err != nil:
			w.WriteHeader(http.StatusInternalServerError)
			json, _ := json.Marshal(ApiError{Error: "internal server error"})
			w.Write(json)
			log.Printf("json.Marshal error: %s", err)
		case result.Error == "not found":
			w.WriteHeader(http.StatusNotFound)
			json, _ := json.Marshal(ApiError{Error: "not found"})
			w.Write(json)
		case !result.Success:
			w.WriteHeader(http.StatusInternalServerError)
			json, _ := json.Marshal(ApiError{Error: result.Error})
			w.Write(json)
		default:
			w.WriteHeader(http.StatusOK)
			w.Write(jsonResp)
		}
	} else {
		result := GetRepoInfo(username, reponame)
		jsonResp, err := json.Marshal(result)
		switch {
		case err != nil:
			w.WriteHeader(http.StatusInternalServerError)
			json, _ := json.Marshal(ApiError{Error: "internal server error"})
			w.Write(json)
			log.Printf("json.Marshal error: %s", err)
		case result.Error == "not found":
			w.WriteHeader(http.StatusNotFound)
			json, _ := json.Marshal(ApiError{Error: "not found"})
			w.Write(json)
		case !result.Success:
			w.WriteHeader(http.StatusInternalServerError)
			json, _ := json.Marshal(ApiError{Error: result.Error})
			w.Write(json)
		default:
			w.WriteHeader(http.StatusOK)
			w.Write(jsonResp)
		}
	}

}
