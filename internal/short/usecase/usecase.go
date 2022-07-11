package usecase

import (
	"errors"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/short/model"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/short/repo"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/pkg/random"
	"github.com/sirupsen/logrus"
)

type shortUsecase struct {
	repo model.ShortRepository
}

func New(repository model.ShortRepository, log *logrus.Logger) model.ShortUsecase {
	return &shortUsecase{
		repo: repository,
	}
}

func (s *shortUsecase) CreateNewShort(location string, userID string) (*model.Short, error) {
	short := model.Short{
		Key:      random.String(8),
		Location: location,
		UserID:   userID,
	}
	err := s.repo.Insert(&short)
	if err != nil {
		if errors.Is(err, repo.ErrDuplicate) {
			return nil, ErrDuplicate
		}
		return nil, err
	}
	return &short, nil
}

func (s *shortUsecase) FindByKey(key string) (*model.Short, error) {

	short, err := s.repo.GetOneByKey(key)
	if err != nil {
		return nil, err
	}

	return short, nil
}

func (s *shortUsecase) FindByLocation(location string) (*model.Short, error) {

	short, err := s.repo.GetOneByKey(location)
	if err != nil {
		return nil, err
	}

	return short, nil
}

func (s *shortUsecase) GetUserShorts(userID string) ([]*model.Short, error) {
	shorts, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	return shorts, nil
}

func (s *shortUsecase) Save(short *model.Short) error {

	err := s.repo.Insert(short)
	if err != nil {
		return err
	}

	return nil
}

func (s *shortUsecase) Delete(key ...string) error {
	return nil
}

type makeShortsBatchDTO struct {
	Location      string
	CorrelationID string
	UserID        string
}

func (s *shortUsecase) CreateNewShortBatch(items ...makeShortsBatchDTO) (map[string]*model.Short, error) {

	shorts := make(map[string]*model.Short, len(items))

	for _, item := range items {
		short, err := s.CreateNewShort(item.Location, item.UserID)
		if err != nil {
			continue
		}
		shorts[item.CorrelationID] = short
	}
	return shorts, nil
}
