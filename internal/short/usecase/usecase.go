package usecase

import (
	"log"
	"yandex-practicum-go-shortener/internal/short/model"
)

type shortUsecase struct {
	repo model.ShortRepository
}

func New(repository model.ShortRepository) model.ShortUsecase {
	return &shortUsecase{
		repo: repository,
	}
}

func (s *shortUsecase) CreateNewShort(location string, userID string) (*model.Short, error) {
	log.Println("make new", location, userID)
	return &model.Short{}, nil
}

func (s *shortUsecase) FindByKey(key string) (*model.Short, error) {
	return &model.Short{}, nil
}

func (s *shortUsecase) FindByLocation(location string) (*model.Short, error) {
	return &model.Short{}, nil
}

func (s *shortUsecase) GetUserShorts(userID string) ([]*model.Short, error) {
	return nil, nil
}

func (s *shortUsecase) Save(*model.Short) error {
	return nil
}

func (s *shortUsecase) Delete(key ...string) error {
	return nil
}
