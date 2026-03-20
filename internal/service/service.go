package service

import (
	"errors"
	"strings"
	"unicode"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func Convert(input string) (string, error) {
	if input == "" {
		return "", errors.New("пустая строка")
	}

	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", errors.New("строка состоит только из пробелов")
	}

	isText := false
	for _, r := range trimmed {
		if unicode.IsLetter(r) {
			isText = true
			break
		}
	}

	var result string
	if isText {
		result = morse.ToMorse(trimmed)
	} else {
		result = morse.ToText(trimmed)
	}

	if result == "" && trimmed != "" {
		return "", errors.New("не удалось выполнить конвертацию: результат пуст")
	}

	return result, nil
}
