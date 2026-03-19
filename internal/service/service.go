package service

import (
	"errors"
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func Convert(input string) (string, error) {
	if input == "" {
		return "", errors.New("пустая строка")
	}
	trimmed := strings.TrimSpace(input)

	isMorse := true
	for _, r := range trimmed {
		if r != '.' && r != '-' && r != ' ' && r != '\n' && r != '\r' {
			isMorse = false
			break
		}
	}
	if isMorse {
		result := morse.ToText(trimmed)
		if result == "" && trimmed != "" {

		}
		return result, nil
	} else {
		return morse.ToMorse(trimmed), nil
	}
}
