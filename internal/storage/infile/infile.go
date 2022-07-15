//Used for file storage. Making file and write url per line
package infile

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"

	"yandex-practicum-go-shortener/internal/storage"
)

type store struct {
	fileMutex sync.Mutex
	file      *os.File
	encoder   *json.Encoder
	decoder   *json.Decoder
}

var _ storage.Storage = (*store)(nil)

func New(filepath string) (storage.Storage, error) {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return &store{}, err
	}

	return &store{
		file:    file,
		decoder: json.NewDecoder(file),
		encoder: json.NewEncoder(file),
	}, nil
}

func (s *store) GetByShort(key string) (storage.URLRecord, error) {
	s.fileMutex.Lock()
	defer s.fileMutex.Unlock()

	_, err := s.file.Seek(0, io.SeekStart)
	if err != nil {
		return storage.URLRecord{}, err
	}
	for {
		var rec storage.URLRecord

		err := s.decoder.Decode(&rec)
		if err != nil {
			break
		}

		if rec.Short == key {
			return rec, nil
		}
	}
	return storage.URLRecord{}, errors.New("errNotFound")
}

func (s *store) SelectByUID(uid string) ([]storage.URLRecord, error) {
	s.fileMutex.Lock()
	defer s.fileMutex.Unlock()

	var result []storage.URLRecord

	_, err := s.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	for {
		var rec storage.URLRecord

		err := s.decoder.Decode(&rec)
		if err != nil {
			break
		}

		if rec.UserID == uid {
			result = append(result, rec)
		}
	}
	return result, err
}

func (s *store) GetByOriginal(original string) (storage.URLRecord, error) {
	s.fileMutex.Lock()
	defer s.fileMutex.Unlock()

	_, err := s.file.Seek(0, io.SeekStart)
	if err != nil {
		return storage.URLRecord{}, err
	}
	for {
		var rec storage.URLRecord

		err := s.decoder.Decode(&rec)
		if err != nil {
			break
		}

		if rec.Original == original {
			return rec, nil
		}
	}
	return storage.URLRecord{}, errors.New("notFound")
}

func (s *store) Save(short, original, userID string) error {
	s.fileMutex.Lock()
	defer s.fileMutex.Unlock()

	_, err := s.file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	rec := storage.URLRecord{
		Short:    short,
		Original: original,
		UserID:   userID,
	}
	if err := s.encoder.Encode(&rec); err != nil {
		return err
	}
	return nil
}

func (s *store) Ping() error {
	return nil
}

func (s *store) Delete(uid string, keys ...string) {}
