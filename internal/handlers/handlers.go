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
		log.Printf("Ошибка копирования файла в ответ: %v", err)
	}
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)

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

	content := string(data)
	if strings.TrimSpace(content) == "" {
		http.Error(w, "Файл пуст", http.StatusBadRequest)
		return
	}

	resultStr, err := service.Convert(content)
	if err != nil {
		http.Error(w, "Ошибка конвертации", http.StatusInternalServerError)
		log.Printf("Ошибка конвертации: %v", err)
		return
	}

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Printf("Ошибка создания директории %s: %v", uploadDir, err)
	}

	timestamp := time.Now().UTC().Format("2006-01-02_15-04-05")
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".txt"
	}

	newFileName := fmt.Sprintf("%sresult_%s%s", uploadDir, timestamp, ext)

	outFile, err := os.Create(newFileName)
	if err != nil {
		log.Printf("Ошибка создания файла %s: %v", newFileName, err)
		http.Error(w, "Ошибка сохранения результата", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = outFile.WriteString(resultStr)
	if err != nil {
		log.Printf("Ошибка записи в файл %s: %v", newFileName, err)
		http.Error(w, "Ошибка записи результата", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, resultStr)

}
