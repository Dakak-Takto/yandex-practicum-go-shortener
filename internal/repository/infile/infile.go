//Used for file storage. Making file and write url per line
package infile

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/entity"
)

type store struct {
	fileMutex sync.Mutex
	file      *os.File
	encoder   *json.Encoder
	decoder   *json.Decoder
}

var _ entity.URLRepository = (*store)(nil)

func New(filepath string) (entity.URLRepository, error) {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return &store{}, err
	}

	return &store{
		file:    file,
		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
	}, nil
}

func (s *store) GetByShort(ctx context.Context, short string) (*entity.URL, error) {
	s.fileMutex.Lock()
	defer s.fileMutex.Unlock()

	_, err := s.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	for {
		var url entity.URL

		if err := s.decoder.Decode(&url); err != nil {
			break
		}

		if url.Short == short {
			return &url, nil
		}
	}
	return nil, entity.ErrNotFound
}

func (s *store) SelectByUserID(ctx context.Context, userID string) (urls []*entity.URL, err error) {
	s.fileMutex.Lock()
	defer s.fileMutex.Unlock()

	_, err = s.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	for {
		var url entity.URL

		if err := s.decoder.Decode(&url); err != nil {
			break
		}

		if url.UserID == userID {
			urls = append(urls, &url)
		}
	}

	if len(urls) == 0 {
		return nil, entity.ErrNotFound
	}

	return urls, nil
}

func (s *store) GetByOriginal(ctx context.Context, original string) (*entity.URL, error) {
	s.fileMutex.Lock()
	defer s.fileMutex.Unlock()

	_, err := s.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	for {
		var url entity.URL

		if err := s.decoder.Decode(&url); err != nil {
			break
		}

		if url.Original == original {
			return &url, nil
		}
	}
	return nil, entity.ErrNotFound
}

func (s *store) Save(ctx context.Context, url *entity.URL) error {
	s.fileMutex.Lock()
	defer s.fileMutex.Unlock()

	_, err := s.file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	if err := s.encoder.Encode(url); err != nil {
		return err
	}

	return nil
}

func (s *store) Delete(uid string, keys ...string) {}
