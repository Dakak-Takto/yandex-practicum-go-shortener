package inmem

import (
	"context"
	"sync"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/entity"
)

type store struct {
	dataMutex sync.Mutex
	data      []*entity.URL
}

var _ entity.URLRepository = (*store)(nil)

func New() (entity.URLRepository, error) {

	return &store{}, nil
}

func (s *store) GetByShort(ctx context.Context, short string) (*entity.URL, error) {
	for _, url := range s.data {
		if url.Short == short {
			return url, nil
		}
	}
	return nil, entity.ErrNotFound
}

func (s *store) GetByOriginal(ctx context.Context, original string) (*entity.URL, error) {
	for _, entity := range s.data {
		if entity.Original == original {
			return entity, nil
		}
	}
	return nil, entity.ErrNotFound
}

func (s *store) SelectByUserID(ctx context.Context, userID string) (urls []*entity.URL, err error) {
	for _, entity := range s.data {
		if entity.UserID == userID {
			urls = append(urls, entity)
		}
	}
	if len(urls) == 0 {
		return nil, entity.ErrNotFound
	}

	return urls, nil
}

func (s *store) Save(ctx context.Context, url *entity.URL) error {
	s.data = append(s.data, url)
	return nil
}

func (s *store) Delete(uid string, shorts ...string) {}
