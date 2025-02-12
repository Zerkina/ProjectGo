package shortener

import (
	"crypto/rand"
	"encoding/base64"
	// "errors"
	// "sync"
)

// URLStore хранит соответствия между короткими и длинными URL
type URLStore struct {
	urlMap map[string]string
}

// NewURLStore создает новый экземпляр URLStore
func NewURLStore() *URLStore {
	return &URLStore{
		urlMap: make(map[string]string),
	}
}

type Shortener interface {
	ShortenURL(originalURL string) string
	ExpandURL(shortID string) (string, bool)
}

// ExpandURL возвращает оригинальный URL по короткому ID
func (s *URLStore) ExpandURL(shortID string) (string, bool) {
	originalURL, ok := s.urlMap[shortID]
	return originalURL, ok
}

// ShortenURL сокращает URL и сохраняет, возвращая только shortID
func (s *URLStore) ShortenURL(originalURL string) string {
	// Генерируем уникальный shortID
	shortID := generateShortID()

	s.urlMap[shortID] = originalURL // Сохраняем originalURL по ключу shortID
	return shortID                  // Возвращаем только shortID
}

// generateShortID генерирует случайный shortID (base64 encoded)
func generateShortID() string {
	// Генерируем 12 байт случайных данных
	bytes := make([]byte, 12)
	if _, err := rand.Read(bytes); err != nil {
		panic(err) // В production коде нужно обрабатывать эту ошибку корректно
	}

	// Кодируем в base64 (URL-safe)
	shortID := base64.URLEncoding.EncodeToString(bytes)

	// Обрезаем до нужной длины (например, 8 символов)
	return shortID[:8]
}
