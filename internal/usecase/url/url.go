package url

import (
	"context"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/entity"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/pkg/random"
	"github.com/sirupsen/logrus"
)

type usecase struct {
	repo entity.URLRepository
	log  *logrus.Logger
}

func New(repo entity.URLRepository, log *logrus.Logger) entity.URLUsecase {
	return &usecase{
		repo: repo,
		log:  log,
	}
}

func (u *usecase) Create(original string, userID string) (*entity.URL, error) {
	url := entity.URL{
		Short:    random.String(8),
		Original: original,
		UserID:   userID,
	}

	if err := u.repo.Save(context.Background(), &url); err != nil {
		return nil, err
	}

	return &url, nil
}

func (u *usecase) Delete(userID string, shorts ...string) {
	u.repo.Delete(userID, shorts...)
}

func (u *usecase) GetByOriginal(original string) (*entity.URL, error) {
	url, err := u.repo.GetByOriginal(context.Background(), original)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (u *usecase) GetByShort(short string) (*entity.URL, error) {
	url, err := u.repo.GetByShort(context.Background(), short)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (u *usecase) UserURLs(userID string) ([]*entity.URL, error) {
	urls, err := u.repo.SelectByUserID(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	return urls, nil
}
