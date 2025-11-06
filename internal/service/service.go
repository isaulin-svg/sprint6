package service

import (
	"strings"

	"github.com/isaulin-svg/sprint6/pkg/morse"
)

// ConvertAutomatically — определяет тип содержимого и конвертирует его.
func ConvertAutomatically(input string) (string, error) {
	trimmed := strings.TrimSpace(input)

	// если строка состоит только из точек, тире и пробелов — это Морзе
	if isMorse(trimmed) {
		return morse.ToText(trimmed), nil
	}

	// если обычный текст — конвертируем в Морзе
	return morse.ToMorse(trimmed), nil
}

// Проверяем, похожа ли строка на Морзе
func isMorse(s string) bool {
	for _, r := range s {
		if r != '.' && r != '-' && r != ' ' && r != '\n' && r != '\r' {
			return false
		}
	}
	return true
}
