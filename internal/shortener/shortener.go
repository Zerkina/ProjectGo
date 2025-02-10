package shortener

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
)

// URLStore хранит соответствия между короткими и длинными URL
type URLStore struct {
	mu     sync.RWMutex
	urlMap map[string]string
}

// NewURLStore создает новый экземпляр URLStore
func NewURLStore() *URLStore {
	return &URLStore{
		urlMap: make(map[string]string),
	}
}

// GetOriginalURL возвращает оригинальный URL по короткому ID
func (s *URLStore) GetOriginalURL(id string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	originalURL, ok := s.urlMap[id]
	if !ok {
		return "", errors.New("URL not found")
	}
	return originalURL, nil
}

// ShortenURL сокращает URL и сохраняет, возвращая только shortID
func (s *URLStore) ShortenURL(originalURL string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Генерируем уникальный shortID
	shortID := generateShortID()

	// Проверяем, не существует ли уже такой shortID. Если существует, генерируем новый.
	for _, ok := s.urlMap[shortID]; ok; {
		shortID = generateShortID()
	}
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
