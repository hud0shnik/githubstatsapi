package convert

import (
	"strconv"
	"strings"
)

// ToInt переводит string в int
func ToInt(s string) int {

	// Удаление запятых из числа
	s = strings.ReplaceAll(s, ",", "")

	i, _ := strconv.Atoi(s)

	return i

}

// ToBool переводит string в bool
func ToBool(s string) bool {

	f, _ := strconv.ParseBool(s)

	return f

}
