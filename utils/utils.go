package utils

import (
	"strconv"
	"strings"
)

// Функция поиска. Возвращает искомое значение и индекс последнего символа
func FindWithIndex(str, subStr, stopChar string, start int) (string, int) {

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
func Find(str, subStr, stopChar string) string {

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
func ToInt(s string) int {

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return i

}

// Функция перевода строки в bool
func ToBool(s string) bool {

	b, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}

	return b

}
