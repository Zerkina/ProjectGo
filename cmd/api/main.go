package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// Функция для генерации короткого URL (заглушка)
func shortenURL(originalURL string) string {
	// Пока просто возвращаем обрезанную версию URL для примера
	if len(originalURL) > 8 {
		return originalURL[0:8]
	}
	return originalURL
}

func mainPage(res http.ResponseWriter, req *http.Request) {
	// 1. Проверяем метод запроса - он должен быть POST
	if req.Method != http.MethodPost {
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

	// 3. Читаем тело запроса (это должен быть URL)
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

	//5.  Сокращаем URL
	shortURL := shortenURL(originalURL)

	// 6.  Устанавливаем Content-Type для ответа как text/plain
	res.Header().Set("Content-Type", "text/plain")

	//7.  Отправляем код ответа 201 Created и сокращённый URL
	res.WriteHeader(http.StatusCreated)
	fmt.Fprint(res, shortURL)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainPage)

	log.Println("Server listening on port 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
