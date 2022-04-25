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

func (s *store) GetByShort(key string) (storage.URLRecord, error) {
	for _, entity := range s.data {
		if entity.Short == key {
			return entity, nil
		}
	}
	return storage.URLRecord{}, errors.New("notFoundError")
}

func (s *store) GetByOriginal(original string) (storage.URLRecord, error) {
	for _, entity := range s.data {
		if entity.Original == original {
			return entity, nil
		}
	}
	return storage.URLRecord{}, errors.New("notFoundError")
}

func (s *store) GetByUID(uid string) ([]storage.URLRecord, error) {
	var result []storage.URLRecord
	for _, entity := range s.data {
		if entity.UserID == uid {
			result = append(result, entity)
		}
	}
	return result, nil
}

func (s *store) Save(short, original, userID string) error {
	s.data = append(s.data, storage.URLRecord{
		Short:    short,
		Original: original,
		UserID:   userID,
	})
	return nil
}

func (s *store) IsExist(key string) bool {
	_, err := s.GetByShort(key)
	return err == nil
}

func (s *store) Lock() {
	s.dataMutex.Lock()
}

func (s *store) Unlock() {
	s.dataMutex.Unlock()
}

func (s *store) Ping() error {
	return nil
}
