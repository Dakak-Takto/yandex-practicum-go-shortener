package inmem

import (
	"errors"
	"sync"
	"yandex-practicum-go-shortener/internal/storage"
)

type store struct {
	dataMutex sync.Mutex
	data      []storage.URLRecord
}

var _ storage.Storage = (*store)(nil)

func New() (storage.Storage, error) {

	return &store{}, nil
}

func (s *store) First(key string) (storage.URLRecord, error) {
	for _, entity := range s.data {
		if entity.Short == key {
			return entity, nil
		}
	}
	return storage.URLRecord{}, errors.New("notFoundError")
}

func (s *store) Get(key string) []storage.URLRecord {
	var result []storage.URLRecord
	for _, entity := range s.data {
		if entity.Short == key {
			result = append(result, entity)
		}
	}
	return result
}

func (s *store) Save(short, original, userID string) {
	s.data = append(s.data, storage.URLRecord{
		Short:    short,
		Original: original,
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
