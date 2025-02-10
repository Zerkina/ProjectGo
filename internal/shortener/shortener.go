package shortener

import (
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

// ShortenURL сокращает URL и сохраняет
func (s *URLStore) ShortenURL(originalURL string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	shortURL := shortenString(originalURL) // Используем функцию сокращения

	// Проверяем, не существует ли уже такой короткий URL. Если существует, добавляем случайный элемент "х"
	for _, ok := s.urlMap[shortURL]; ok; {
		shortURL = shortenString(originalURL + "x")
	}
	s.urlMap[shortURL] = originalURL
	return shortURL
}

func shortenString(originalURL string) string {
	if len(originalURL) > 8 {
		return originalURL[:8]
	}
	return originalURL
}
