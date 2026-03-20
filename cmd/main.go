package main

import (
	"log"
	"os"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "[MorseConverter] ", log.Ldate|log.Ltime|log.Lshortfile)

	srv := server.NewServer(logger)

	logger.Println("Запуск сервера на порту 8080...")

	if err := srv.Svr.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}

}
