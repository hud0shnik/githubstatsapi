package api2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hud0shnik/githubstatsapi/utils"
	"github.com/sirupsen/logrus"
)

// Структура для хранения полной информации о пользователе
type userInfo struct {
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
type userInfoString struct {
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

// Структура ошибки
type apiError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// Функция получения информации о пользователе в формате строк
func getUserInfoString(username string) (userInfoString, int, error) {

	// Формирование и исполнение запроса
	resp, err := http.Get("https://github.com/" + username)
	if err != nil {
		return userInfoString{}, http.StatusInternalServerError, fmt.Errorf("in http.Get: %w", err)
	}
	defer resp.Body.Close()

	// Проверка статускода
	if resp.StatusCode != 200 {
		return userInfoString{}, resp.StatusCode,
			fmt.Errorf(resp.Status)
	}

	// Запись респонса
	body, _ := ioutil.ReadAll(resp.Body)

	// HTML полученной страницы в формате string
	pageStr := string(body)[55000:]

	// Запись html в файл для тестирования
	/*if err := os.WriteFile("sample.html", []byte(pageStr), 0666); err != nil {
		logrus.Fatal(err)
	}*/

	// Проверка на страницу пользователя
	if !strings.Contains(pageStr, "p-nickname vcard-username d-block") {
		return userInfoString{}, http.StatusNotFound, fmt.Errorf("not found")
	}

	// Проверка на скрытие коммитов
	if strings.Contains(pageStr, "'s activity is private</h4>") {
		return userInfoString{}, http.StatusForbidden, fmt.Errorf("activity is private")
	}

	// Структура, которую будет возвращать функция
	result := userInfoString{
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

	return result, http.StatusOK, nil

}

// Функция получения информации о пользователе в формате строк
func getUserInfo(username string) (userInfo, int, error) {

	// Получение текстовой версии статистики
	resultStr, statusCode, err := getUserInfoString(username)
	if err != nil {
		return userInfo{}, statusCode, err
	}

	// Форматирование
	return userInfo{
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
	}, http.StatusOK, nil

}

// Роут "/user"
func User(w http.ResponseWriter, r *http.Request) {

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

	// Проверка на тип
	if r.URL.Query().Get("type") == "string" {

		// Получение статистики
		result, statusCode, err := getUserInfoString(id)
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

		// Получение статистики
		result, statusCode, err := getUserInfo(id)
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
