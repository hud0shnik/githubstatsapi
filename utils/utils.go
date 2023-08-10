package utils

import (
	"strconv"
	"strings"
)

// FindWithIndex производит поиск substr в s[start:] и возвращает строку от конца substr до stopChar
func FindWithIndex(s, substr, stopChar string, start int) (string, int) {

	// Обрезка левой границы поиска
	s = s[start:]

	// Проверка на существование нужной строки
	if strings.Contains(s, substr) {

		// Поиск индекса начала нужной строки
		left := strings.Index(s, substr) + len(substr)

		// Поиск правой границы
		right := left + strings.Index(s[left:], stopChar)

		// Обрезка и вывод результата
		return s[left:right], right + start
	}

	return "", 0

}

// Find работает как FindWithIndex, но не возвращает индекс
func Find(s, substr, stopChar string) string {

	// Проверка на существование нужной строки
	if strings.Contains(s, substr) {

		// Обрезка левой части
		s = s[strings.Index(s, substr)+len(substr):]

		// Обрезка правой части и вывод результата
		return s[:strings.Index(s, stopChar)]
	}

	return ""

}

// ToInt переводит string в int
func ToInt(s string) int {

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return i

}

// ToBool переводит string в bool
func ToBool(s string) bool {

	f, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}

	return f

}
