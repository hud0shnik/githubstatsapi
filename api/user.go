package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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

// Функция поиска. Возвращает искомое значение и индекс последнего символа
func findWithIndex(str, subStr, stopChar string, start int) (string, int) {

	// Обрезка левой границы поиска
	str = str[start:]

	// Проверка на существование нужной строки
	if strings.Contains(str, subStr) {

		// Поиск индекса начала нужной строки
		left := strings.Index(str, subStr) + len(subStr)

		// Поиск правой границы
		right := left + strings.Index(str[left:], stopChar)

		// Обрезка и вывод результата
		return str[left:right], right + start
	}

	return "", 0

}

// Облегчённая функция поиска. Возвращает только искомое значение
func find(str, subStr, stopChar string) string {

	// Проверка на существование нужной строки
	if strings.Contains(str, subStr) {

		// Обрезка левой части
		str = str[strings.Index(str, subStr)+len(subStr):]

		// Обрезка правой части и вывод результата
		return str[:strings.Index(str, stopChar)]
	}

	return ""

}

// Функция перевода строки в число
func toInt(s string) int {

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return i

}

// Функция перевода строки в bool
func toBool(s string) bool {

	b, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}

	return b

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

	// Запись респонса
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

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
	result.Repositories, left = findWithIndex(pageStr, "Repositories\n    <span title=\"", "\"", left)

	// Пакеты
	result.Packages, left = findWithIndex(pageStr, "Packages\n      <span title=\"", "\"", left)

	// Поставленные звезды
	result.Stars, left = findWithIndex(pageStr, "Stars\n    <span title=\"", "\"", left)

	// Ссылка на аватар
	result.Avatar, left = findWithIndex(pageStr, "<img style=\"height:auto;\" alt=\"Avatar\" src=\"", "\"", left)

	// Статус
	result.Status, left = findWithIndex(pageStr, "status-message-wrapper f6 color-fg-default no-wrap \" >\n        <div>", "</div>", left)

	// Имя пользователя
	result.Name, left = findWithIndex(pageStr, "\"name\">\n          ", "\n", left)

	// Подписчики
	result.Followers, left = findWithIndex(pageStr, "<span class=\"text-bold color-fg-default\">", "<", left)

	// Подписки
	result.Following, left = findWithIndex(pageStr, "<span class=\"text-bold color-fg-default\">", "<", left)

	// Контрибуции за год
	result.Contributions, _ = findWithIndex(pageStr, "<h2 class=\"f4 text-normal mb-2\">\n      ", "\n", left)

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
		Followers:     toInt(resultStr.Followers),
		Following:     toInt(resultStr.Following),
		Repositories:  toInt(resultStr.Repositories),
		Packages:      toInt(resultStr.Packages),
		Stars:         toInt(resultStr.Stars),
		Contributions: toInt(resultStr.Contributions),
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
