package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/Zerkina/url-shortener/internal/handlers"
)

func main() {
	// Инициализируем обработчик
	h := handlers.NewHandler()

	// Регистрируем обработчики для маршрутов
	http.HandleFunc("/", h.MainPage) // POST запросы на / обрабатывает MainPage (создание короткой ссылки)

	// Обработчик для редиректа. Здесь важна правильная обработка пути.
	http.HandleFunc("/redirect/", func(w http.ResponseWriter, r *http.Request) {
		// Обрезаем "/redirect/" из пути, чтобы получить ID.
		id := strings.TrimPrefix(r.URL.Path, "/redirect/")

		// Передаем обрезанный ID в RedirectHandler
		r.URL.Path = "/" + id // Изменяем путь запроса, чтобы RedirectHandler получил ID
		h.RedirectHandler(w, r)
	})

	// Запускаем сервер
	log.Fatal(http.ListenAndServe(":8080", nil))
}
