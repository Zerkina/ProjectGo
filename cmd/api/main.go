package main

import (
	"log"
	"net/http"

	// "strings"

	"github.com/Zerkina/url-shortener/internal/handlers"
)

func main() {
	// 1. Создаем URLStore

	// Инициализируем обработчикurlStore
	h := handlers.NewHandler()

	// Регистрируем обработчики для маршрутов
	http.HandleFunc("/", h.MainPage) // POST запросы на / обрабатывает MainPage (создание короткой ссылки)
	http.HandleFunc("/{id}", h.Redirect)

	// Запускаем сервер
	log.Fatal(http.ListenAndServe(":8080", nil))
}
