package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/Zerkina/url-shortener/internal/shortener"
)

type Handler struct {
	shortener *shortener.URLStore
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
	fmt.Fprint(res, shortID) // Возвращаем *только* shortID

	// 6. Устанавливаем Content-Type для ответа как text/plain
	res.Header().Set("Content-Type", "text/plain")

	// 7. Отправляем код ответа 201 Created и сокращённый URL
	res.WriteHeader(http.StatusCreated)
	fmt.Fprint(res, shortID)
}

// RedirectHandler обрабатывает GET запросы для редиректа
func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID из URL
	id := strings.TrimPrefix(r.URL.Path, "/")

	// Получаем оригинальный URL из shortener (используем хранилище)
	originalURL, err := h.shortener.GetOriginalURL(id)
	if err != nil {
		log.Println("Error getting original URL:", err)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// Выполняем редирект
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
