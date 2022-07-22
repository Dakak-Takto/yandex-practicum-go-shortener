//Package infile Used for file storage. Making file and write url per line
package infile

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/storage"
)

type store struct {
	fileMutex sync.Mutex
	file      *os.File
}

var _ storage.Storage = (*store)(nil)

func New(filepath string) (storage.Storage, error) {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return &store{}, err
	}

	return &store{
		file: file,
	}, nil
}

func (s *store) GetByShort(key string) (storage.URLRecord, error) {
	s.fileMutex.Lock()
	defer s.fileMutex.Unlock()

	s.file.Seek(0, io.SeekStart)

	decoder := json.NewDecoder(s.file)

	for {
		var rec storage.URLRecord

		if err := decoder.Decode(&rec); err != nil {
			fmt.Printf("error decode: %s", err)
			break
		}

		fmt.Printf("found item: %v\n", rec)

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

	s.file.Seek(0, io.SeekStart)
	decoder := json.NewDecoder(s.file)

	for {
		var rec storage.URLRecord

		err := decoder.Decode(&rec)
		if err != nil {
			break
		}

		if rec.UserID == uid {
			result = append(result, rec)
		}
	}
	return result, nil
}

func (s *store) GetByOriginal(original string) (storage.URLRecord, error) {
	s.fileMutex.Lock()
	defer s.fileMutex.Unlock()

	s.file.Seek(0, io.SeekStart)
	decoder := json.NewDecoder(s.file)

	for {
		var rec storage.URLRecord

		err := decoder.Decode(&rec)
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

	s.file.Seek(0, io.SeekEnd)

	rec := storage.URLRecord{
		Short:    short,
		Original: original,
		UserID:   userID,
	}
	if err := json.NewEncoder(s.file).Encode(rec); err != nil {
		return err
	}
	return nil
}

func (s *store) Ping() error {
	return nil
}

func (s *store) Delete(uid string, keys ...string) {}
