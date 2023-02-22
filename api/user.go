package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Структура для хранения полной информации о пользователе
type UserInfo struct {
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

// Функция получения информации о пользователе
func GetUserInfo(username string) UserInfo {

	// Формирование и исполнение запроса
	resp, err := http.Get("https://github.com/" + username)
	if err != nil {
		return UserInfo{
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
		return UserInfo{
			Error: "user not found",
		}
	}

	// Проверка на скрытие коммитов
	if strings.Contains(pageStr, "'s activity is private</h4>") {
		return UserInfo{
			Error: username + "'s activity is private",
		}
	}

	// Структура, которую будет возвращать функция
	result := UserInfo{
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

// Роут "/user"
func User(w http.ResponseWriter, r *http.Request) {

	// Получение параметра id из реквеста
	id := r.URL.Query().Get("id")

	// Если параметра нет, отправка ошибки
	if id == "" {
		http.NotFound(w, r)
		return
	}

	// Передача в заголовок респонса типа данных
	w.Header().Set("Content-Type", "application/json")

	// Форматирование структуры в json и отправка пользователю
	jsonResp, err := json.Marshal(GetUserInfo(id))
	if err != nil {
		log.Printf("json.Marshal error: %s", err)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	}
}
