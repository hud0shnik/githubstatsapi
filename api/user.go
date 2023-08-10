package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/hud0shnik/githubstatsapi/utils"
)

// Структура для хранения полной информации о пользователе
type UserInfo struct {
	Success       bool   `json:"success"`
	Error         string `json:"error"`
	Username      string `json:"username"`
	Name          string `json:"name"`
	Followers     int    `json:"followers"`
	Following     int    `json:"following"`
	Repositories  int    `json:"repositories"`
	Packages      int    `json:"packages"`
	Stars         int    `json:"stars"`
	Contributions int    `json:"contributions"`
	Status        string `json:"status"`
	Avatar        string `json:"avatar"`
}

// Структура для парсинга полной информации о пользователе
type UserInfoString struct {
	Success       bool   `json:"success"`
	Error         string `json:"error"`
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

// Функция получения информации о пользователе в формате строк
func GetUserInfoString(username string) UserInfoString {

	// Формирование и исполнение запроса
	resp, err := http.Get("https://github.com/" + username)
	if err != nil {
		return UserInfoString{
			Error: "Cant reach github.com",
		}
	}
	defer resp.Body.Close()

	// Запись респонса
	body, _ := io.ReadAll(resp.Body)

	// HTML полученной страницы в формате string
	pageStr := string(body)[55000:]

	// Запись html в файл для тестирования
	/*if err := os.WriteFile("sample.html", []byte(pageStr), 0666); err != nil {
		log.Fatal(err)
	}*/

	// Проверка на страницу пользователя
	if !strings.Contains(pageStr, "p-nickname vcard-username d-block") {
		return UserInfoString{
			Error: "user not found",
		}
	}

	// Проверка на скрытие коммитов
	if strings.Contains(pageStr, "'s activity is private</h4>") {
		return UserInfoString{
			Error: username + "'s activity is private",
		}
	}

	// Структура, которую будет возвращать функция
	result := UserInfoString{
		Success:  true,
		Username: username,
	}

	// Индекс конца последней найденной строки
	left := 0

	// Репозитории
	result.Repositories, left = utils.FindWithIndex(pageStr, "Repositories\n    <span title=\"", "\"", left)

	// Пакеты
	result.Packages, left = utils.FindWithIndex(pageStr, "Packages\n      <span title=\"", "\"", left)

	// Поставленные звезды
	result.Stars, left = utils.FindWithIndex(pageStr, "Stars\n    <span title=\"", "\"", left)

	// Ссылка на аватар
	result.Avatar, left = utils.FindWithIndex(pageStr, "<img style=\"height:auto;\" alt=\"Avatar\" src=\"", "\"", left)

	// Статус
	result.Status, left = utils.FindWithIndex(pageStr, "status-message-wrapper f6 color-fg-default no-wrap \" >\n        <div>", "</div>", left)

	// Имя пользователя
	result.Name, left = utils.FindWithIndex(pageStr, "\"name\">\n          ", "\n", left)

	// Подписчики
	result.Followers, left = utils.FindWithIndex(pageStr, "<span class=\"text-bold color-fg-default\">", "<", left)

	// Подписки
	result.Following, left = utils.FindWithIndex(pageStr, "<span class=\"text-bold color-fg-default\">", "<", left)

	// Контрибуции за год
	result.Contributions, _ = utils.FindWithIndex(pageStr, "<h2 class=\"f4 text-normal mb-2\">\n      ", "\n", left)

	return result

}

// Функция получения информации о пользователе в формате строк
func GetUserInfo(username string) UserInfo {

	// Получение текстовой версии статистики
	resultStr := GetUserInfoString(username)

	// Проверка на ошибки при парсинге
	if !resultStr.Success {
		return UserInfo{
			Success: false,
			Error:   resultStr.Error,
		}
	}

	return UserInfo{
		Success:       resultStr.Success,
		Error:         resultStr.Error,
		Username:      username,
		Name:          resultStr.Name,
		Followers:     utils.ToInt(resultStr.Followers),
		Following:     utils.ToInt(resultStr.Following),
		Repositories:  utils.ToInt(resultStr.Repositories),
		Packages:      utils.ToInt(resultStr.Packages),
		Stars:         utils.ToInt(resultStr.Stars),
		Contributions: utils.ToInt(resultStr.Contributions),
		Status:        resultStr.Status,
		Avatar:        resultStr.Avatar,
	}

}

// Роут "/user"
func User(w http.ResponseWriter, r *http.Request) {

	// Передача в заголовок респонса типа данных
	w.Header().Set("Content-Type", "application/json")

	// Получение параметра id из реквеста
	id := r.URL.Query().Get("id")

	// Если параметра нет, отправка ошибки
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(map[string]string{"Error": "Please insert user id"})
		w.Write(json)
		return
	}

	// Проверка на тип, получение статистики, форматирование и отправка
	if r.URL.Query().Get("type") == "string" {
		jsonResp, err := json.Marshal(GetUserInfoString(id))
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
		jsonResp, err := json.Marshal(GetUserInfo(id))
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
