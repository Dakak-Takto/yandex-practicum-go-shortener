package repo

import (
	"os"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/short/model"
)

type fileRepo struct {
	file *os.File
}

func NewFileRepository(filename string) (model.ShortRepository, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return &fileRepo{
		file: file,
	}, nil
}

func (m *fileRepo) GetOneByKey(key string) (*model.Short, error) {
	return nil, nil
}
func (m *fileRepo) GetOneByLocation(location string) (*model.Short, error) {
	return nil, nil
}
func (m *fileRepo) GetByUserID(userID string) ([]*model.Short, error) {
	return nil, nil
}
func (m *fileRepo) Insert(*model.Short) error {
	return nil
}
func (m *fileRepo) Delete(key ...string) error {
	return nil
}
