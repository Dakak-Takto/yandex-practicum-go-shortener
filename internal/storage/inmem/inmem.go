package inmem

import (
	"errors"
	"sync"
	"yandex-practicum-go-shortener/internal/storage"
)

type store struct {
	urlsMutex sync.Mutex
	urls      map[string]string
}

var _ storage.Storage = (*store)(nil)

func New() storage.Storage {

	return &store{
		urls: make(map[string]string),
	}
}

func (s *store) Get(key string) (value string, err error) {
	if value, ok := s.urls[key]; ok {
		return value, nil
	}

	return "", errors.New("not found")
}

func (s *store) Set(key, value string) error {
	s.urls[key] = value

	return nil
}

func (s *store) IsExist(key string) (isExists bool) {
	_, isExists = s.urls[key]

	return isExists
}

func (s *store) Lock() {
	s.urlsMutex.Lock()
}

func (s *store) Unlock() {
	s.urlsMutex.Unlock()
}