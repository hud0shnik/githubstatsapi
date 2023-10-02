package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/hud0shnik/githubstatsapi/utils"
)

// Структура для хранения информации о репозитории
type repoInfo struct {
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
type repoInfoString struct {
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
func GetRepoInfoString(username, reponame string) repoInfoString {

	// Формирование и исполнение запроса
	resp, err := http.Get("https://github.com/" + username + "/" + reponame)
	if err != nil {
		return repoInfoString{
			Error: "Cant reach github.com",
		}
	}
	defer resp.Body.Close()

	// Запись респонса
	body, _ := io.ReadAll(resp.Body)

	// HTML полученной страницы в формате string
	pageStr := string(body)[20000:]

	// Запись html в файл для тестирования
	/*if err := os.WriteFile("sampleRepo.html", []byte(pageStr), 0666); err != nil {
		log.Fatal(err)
	}*/

	// Проверка на репозиторий
	if !strings.Contains(pageStr, "name=\"selected-link\" value=\"repo_source\"") {
		return repoInfoString{
			Error: "repo doesn't exist",
		}
	}

	// Структура, которую будет возвращать функция
	result := repoInfoString{
		Success:  true,
		Username: username,
		Reponame: reponame,
	}

	// Индекс конца последней найденной строки
	left := 0

	// Ветки
	result.Branches, left = utils.FindWithIndex(pageStr, "01-1.5 0z\"></path>\n</svg>\n          <strong>", "<", left)

	// Теги
	result.Tags, left = utils.FindWithIndex(pageStr, "000-2z\"></path>\n</svg>\n        <strong>", "<", left)

	// Коммиты
	result.Commits, left = utils.FindWithIndex(pageStr, "class=\"d-none d-sm-inline\">\n                    <strong>", "<", left)

	// Звезды
	result.Stars, left = utils.FindWithIndex(pageStr, "94v.001z\"></path>\n</svg>\n    <strong>", "<", left)

	// Просмотры
	result.Watching, left = utils.FindWithIndex(pageStr, " 000 4z\"></path>\n</svg>\n    <strong>", "<", left)

	// Форки
	result.Forks, _ = utils.FindWithIndex(pageStr, "5.75.75 0 000 1.5z\"></path>\n</svg>\n    <strong>", "<", left)

	return result

}

// Функция получения информации о репозитории
func GetrepoInfo(username, reponame string) repoInfo {

	// Получение текстовой версии статистики
	resultStr := GetRepoInfoString(username, reponame)

	// Проверка на ошибки при парсинге
	if !resultStr.Success {
		return repoInfo{
			Success: false,
			Error:   resultStr.Error,
		}
	}

	return repoInfo{
		Success:  resultStr.Success,
		Error:    resultStr.Error,
		Username: resultStr.Username,
		Reponame: resultStr.Reponame,
		Commits:  utils.ToInt(resultStr.Commits),
		Branches: utils.ToInt(resultStr.Branches),
		Tags:     utils.ToInt(resultStr.Tags),
		Stars:    utils.ToInt(resultStr.Stars),
		Watching: utils.ToInt(resultStr.Watching),
		Forks:    utils.ToInt(resultStr.Forks),
	}

}

// Роут "/repo"
func Repo(w http.ResponseWriter, r *http.Request) {

	// Установка заголовков
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// Получение параметра имени пользователя и названия репозитория из реквеста
	username := r.URL.Query().Get("username")
	reponame := r.URL.Query().Get("reponame")

	// Если параметра нет, отправка ошибки
	if username == "" || reponame == "" {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(map[string]string{"Error": "Please insert username and reponame"})
		w.Write(json)
		return
	}

	// Проверка на тип, получение статистики, форматирование и отправка
	if r.URL.Query().Get("type") == "string" {
		jsonResp, err := json.Marshal(GetRepoInfoString(username, reponame))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json, _ := json.Marshal(map[string]string{"Error": "Internal Server Error"})
			w.Write(json)
			log.Printf("json.Marshal error: %s", err)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(jsonResp)
		}
	} else {
		jsonResp, err := json.Marshal(GetrepoInfo(username, reponame))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json, _ := json.Marshal(map[string]string{"Error": "Internal Server Error"})
			w.Write(json)
			log.Printf("json.Marshal error: %s", err)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(jsonResp)
		}
	}

}
