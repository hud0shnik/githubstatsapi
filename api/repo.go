package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Структура для хранения информации о репозитории
type RepoInfo struct {
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

// Функция получения информации о репозитории
func GetRepoInfo(username string, reponame string) RepoInfo {

	// Формирование и исполнение запроса
	resp, err := http.Get("https://github.com/" + username + "/" + reponame)
	if err != nil {
		return RepoInfo{
			Error: "http.Get error",
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
		return RepoInfo{
			Error: "repo doesn't exist",
		}
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

// Роут "/repo"
func Repo(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	username := r.URL.Query().Get("username")
	if username == "" {
		http.NotFound(w, r)
		return
	}
	reponame := r.URL.Query().Get("reponame")
	if reponame == "" {
		http.NotFound(w, r)
		return
	}
	resp := GetRepoInfo(username, reponame)

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Print("Error: ", err)
	} else {
		w.Write(jsonResp)
	}
}
