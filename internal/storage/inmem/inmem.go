package inmem

import (
	"errors"
	"sync"

	"yandex-practicum-go-shortener/internal/storage"
)

type store struct {
	dataMutex sync.Mutex
	data      []storage.Entity
}

var _ storage.Storage = (*store)(nil)

func New() (storage.Storage, error) {

	return &store{}, nil
}

func (s *store) First(key string) (storage.Entity, error) {
	for _, entity := range s.data {
		if entity.Key == key {
			return entity, nil
		}
	}
	return storage.Entity{}, errors.New("notFoundError")
}

func (s *store) Get(key string) []storage.Entity {
	var result []storage.Entity
	for _, entity := range s.data {
		if entity.Key == key {
			result = append(result, entity)
		}
	}
	return result
}

func (s *store) Insert(key, value string) {
	s.data = append(s.data, storage.Entity{
		Key:   key,
		Value: value,
	})
}

func (s *store) IsExist(key string) bool {
	_, err := s.First(key)
	return err == nil
}

func (s *store) Lock() {
	s.dataMutex.Lock()
}

func (s *store) Unlock() {
	s.dataMutex.Unlock()
}
