package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

const (
	uploadDir   = "uploads/"
	maxFileSize = 10 << 20
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("index.html")
	if err != nil {
		http.Error(w, "Файл не найден", http.StatusInternalServerError)
		log.Printf("Ошибка открытия index.html: %v", err)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = io.Copy(w, file)
	if err != nil {
		log.Printf("Ошибка копирования файла в ответ: %v, err")
	}
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		http.Error(w, "Ошибка парсинга формы", http.StatusInternalServerError)
		log.Printf("Ошибка парсинга формы: %v", err)
		return
	}

	file, header, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Ошибка получения файла", http.StatusInternalServerError)
		log.Printf("Ошибка получения файла: %v", err)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
		log.Printf("Ошибка чтения файла: %v", err)
		return
	}

	resultStr, err := service.Convert(string(data))
	if err != nil {
		http.Error(w, "Ошибка конвертации", http.StatusInternalServerError)
		log.Printf("Ошибка конвертации: %v", err)
		return
	}

	if resultStr == "" {
		http.Error(w, "Пустой результат конвертации", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Результат конвертации:\n%s", resultStr)

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Printf("Ошибка создания директории для загрузок: %v", err)
	}

	timestamp := time.Now().UTC().String()
	safeTimestamp := strings.ReplaceAll(timestamp, ":", "-")
	safeTimestamp = strings.ReplaceAll(safeTimestamp, " ", "_")

	ext := filepath.Ext(header.Filename)
	newFileName := fmt.Sprintf("%sresult_%s%s", uploadDir, safeTimestamp, ext)

	outFile, err := os.Create(newFileName)
	if err != nil {
		log.Printf("Ошибка создания файла %s: %v", newFileName, err)
		return
	}
	defer outFile.Close()

	_, err = outFile.WriteString(resultStr)
	if err != nil {
		log.Printf("Ошибка записи в файл: %v", err)
	}
}
