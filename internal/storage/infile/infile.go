//Used for file storage. Making file and write url per line
package infile

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
	"sync"

	"yandex-practicum-go-shortener/internal/storage"
)

type store struct {
	fileMutex sync.Mutex
	file      *os.File
	reader    *bufio.Reader
	writer    *bufio.Writer
}

var _ storage.Storage = (*store)(nil)

func New(filepath string) (storage.Storage, error) {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return &store{}, err
	}

	return &store{
		file:   file,
		reader: bufio.NewReader(file),
		writer: bufio.NewWriter(file),
	}, nil
}

func (s *store) GetByShort(key string) (storage.URLRecord, error) {

	_, err := s.file.Seek(0, io.SeekStart)
	if err != nil {
		return storage.URLRecord{}, err
	}
	for {
		b, _, err := s.reader.ReadLine()
		if err != nil {
			break
		}
		record := strings.Split(string(b), ",")
		if key == record[0] {
			return storage.URLRecord{
				Short:    record[0],
				Original: record[1],
			}, nil
		}
	}
	return storage.URLRecord{}, errors.New("errNotFound")
}

func (s *store) SelectByUID(uid string) ([]storage.URLRecord, error) {
	var result []storage.URLRecord

	_, err := s.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	for {
		b, _, err := s.reader.ReadLine()
		if err != nil {
			break
		}
		record := strings.Split(string(b), ",")
		if record[2] == uid {
			result = append(result, storage.URLRecord{
				Short:    record[0],
				Original: record[1],
				UserID:   record[2],
			})
		}
	}
	return result, nil
}

func (s *store) GetByOriginal(original string) (storage.URLRecord, error) {

	_, err := s.file.Seek(0, io.SeekStart)
	if err != nil {
		return storage.URLRecord{}, err
	}
	for {
		b, _, err := s.reader.ReadLine()
		if err != nil {
			break
		}
		record := strings.Split(string(b), ",")
		if record[2] == original {
			return storage.URLRecord{
				Short:    record[0],
				Original: record[1],
				UserID:   record[2],
			}, nil
		}
	}
	return storage.URLRecord{}, errors.New("notFound")
}

func (s *store) Save(short, original, userID string) error {
	_, err := s.file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	record := []string{short, original, userID}

	_, err = s.writer.WriteString(strings.Join(record, ",") + "\n")
	if err != nil {
		return err
	}
	return s.writer.Flush()
}

func (s *store) Lock() {
	s.fileMutex.Lock()
}
func (s *store) Unlock() {
	s.fileMutex.Unlock()
}

func (s *store) Ping() error {
	return nil
}

func (s *store) Delete(uid string, keys ...string) {}
