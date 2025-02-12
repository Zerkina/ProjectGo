package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	// "strings"

	"github.com/Zerkina/url-shortener/internal/shortener"
)

// Handler структура для обработчиков
type Handler struct {
	shortener shortener.Shortener // Тип - интерфейс
}

func NewHandler() *Handler {
	return &Handler{
		shortener: shortener.NewURLStore(),
	}
}

func (h *Handler) MainPage(res http.ResponseWriter, req *http.Request) {
	// 1. Проверяем метод запроса - он должен быть POST
	if req.Method != http.MethodPost {
		log.Println("Method is not POST")
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "Ожидаем POST")
		return
	}

	// 2. Проверяем Content-Type - он должен быть text/plain
	contentType := req.Header.Get("Content-Type")
	if contentType != "text/plain" {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "Ожидаем Content-Type равным text/plain")
		return
	}

	// 3. Читаем тело запроса
	body, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "Не удалось прочитать тело")
		return
	}
	defer req.Body.Close()
	originalURL := string(body)

	// 4. Валидация URL
	_, err = url.ParseRequestURI(originalURL)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(res, "Invalid URL: %v", err)
		return
	}

	// 5. Сокращаем URL
	shortID := h.shortener.ShortenURL(originalURL)
	shortenedURL := fmt.Sprintf("http://localhost:8080/%s", shortID) // Полный URL.
	// fmt.Fprint(res, shortenedURL)

	// 6. Устанавливаем Content-Type для ответа как text/plain
	res.Header().Set("Content-Type", "text/plain")

	// 7. Отправляем код ответа 201 Created и сокращённый URL
	res.WriteHeader(http.StatusCreated)
	fmt.Fprint(res, shortenedURL)
}

// Redirect обработчик для GET запроса /{id} (редирект на оригинальный URL).
func (h *Handler) Redirect(res http.ResponseWriter, req *http.Request) {
	// 1. Проверяем метод запроса
	if req.Method != http.MethodGet {
		log.Println("Method is not GET")
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "Ожидаем GET")
		return
	}

	// 2. Извлекаем shortID из пути URL
	shortID := req.URL.Path[1:] // Получаем ID из пути (обрезаем первый '/')
	log.Printf("Attempting to redirect shortID: %s", shortID)

	// 3. Ищем оригинальный URL по shortID
	originalURL, ok := h.shortener.ExpandURL(shortID)
	if !ok {
		log.Printf("Short URL not found: %s", shortID)
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "Invalid short URL")
		return
	}

	// 4. Устанавливаем заголовок Location и код 307
	res.Header().Set("Location", originalURL)
	res.WriteHeader(http.StatusTemporaryRedirect) // 307 Temporary Redirect
}
