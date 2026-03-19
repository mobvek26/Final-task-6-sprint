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
		log.Printf("Ошибка копирования файл в ответ: %v, err")
	}
}
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
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
		http.Error(w, "Ошибка конертации", http.StatusInternalServerError)
		log.Printf("Ошибка конвертации: %v", err)
		return
	}

	timestamp := time.Now().UTC().String()
	safeTimestamp := strings.ReplaceAll(timestamp, ":", "-")
	safeTimestamp = strings.ReplaceAll(safeTimestamp, " ", "_")

	ext := filepath.Ext(header.Filename)
	newFileName := fmt.Sprintf("result_%s%s", safeTimestamp, ext)

	outFile, err := os.Create(newFileName)
	if err != nil {
		http.Error(w, "Ошибка сщздания файла результата", http.StatusInternalServerError)
		log.Printf("Ошибка создания файла %s: %v", newFileName, err)
		return
	}
	defer outFile.Close()

	_, err = outFile.WriteString(resultStr)
	if err != nil {
		http.Error(w, "Ошибка записи в файл", http.StatusInternalServerError)
		log.Printf("Ошибка записи в файл: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Конвертация успешна\nРезультат сохранен в файл: %s\n\nСодержимое: \n%s", newFileName, resultStr)

}
