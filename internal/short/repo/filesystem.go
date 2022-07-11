package repo

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/short/model"
)

type fileRepo struct {
	mu      sync.Mutex
	file    *os.File
	decoder *json.Decoder
	encoder *json.Encoder
}

func NewFileRepository(filename string) (model.ShortRepository, error) {

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return &fileRepo{
		file:    file,
		decoder: json.NewDecoder(file),
		encoder: json.NewEncoder(file),
	}, nil
}

func (m *fileRepo) GetOneByKey(key string) (*model.Short, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, err := m.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	for {
		var short model.Short
		err := m.decoder.Decode(&short)
		if err != nil {
			break
		}
		if short.Key == key {
			return &short, nil
		}
	}
	return nil, fmt.Errorf("not found")
}
func (m *fileRepo) GetOneByLocation(location string) (*model.Short, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, err := m.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	for {
		var short model.Short
		err := m.decoder.Decode(&short)
		if err != nil {
			break
		}
		if short.Location == location {
			return &short, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (m *fileRepo) GetByUserID(userID string) ([]*model.Short, error) {
	return nil, fmt.Errorf("not implemented")
}
func (m *fileRepo) Insert(short *model.Short) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, err := m.file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	return m.encoder.Encode(short)
}

func (m *fileRepo) Delete(key ...string) error {
	return nil
}
