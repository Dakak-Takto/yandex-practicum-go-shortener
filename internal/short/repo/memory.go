package repo

import (
	"fmt"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/short/model"
)

type memoryRepo struct {
	db map[string]*model.Short
}

func NewMemoryRepository() model.ShortRepository {
	return &memoryRepo{
		db: make(map[string]*model.Short),
	}
}

func (m *memoryRepo) GetOneByKey(key string) (*model.Short, error) {
	short, exist := m.db[key]
	if !exist {
		return nil, fmt.Errorf("not exist key")
	}
	return short, nil
}

func (m *memoryRepo) GetOneByLocation(location string) (*model.Short, error) {
	for _, short := range m.db {
		if short.Location == location {
			return short, nil
		}
	}
	return nil, model.ErrNotFound
}

func (m *memoryRepo) GetByUserID(userID string) ([]*model.Short, error) {

	var shorts []*model.Short

	for _, short := range m.db {
		if short.UserID == userID {
			shorts = append(shorts, short)
		}
	}

	if len(shorts) == 0 {
		return nil, model.ErrNotFound
	}

	return shorts, nil
}

func (m *memoryRepo) Insert(short *model.Short) error {
	m.db[short.Key] = short
	return nil
}

func (m *memoryRepo) Delete(keys ...string) error {

	for _, key := range keys {
		short, err := m.GetOneByKey(key)
		if err != nil {
			continue
		}
		short.Deleted = true
	}
	return nil
}
