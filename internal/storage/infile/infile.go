//Used for file storage. Making file and write url per line
package infile

import (
	"bufio"
	"errors"
	"io"
	"log"
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

func (s *store) First(key string) (storage.URLRecord, error) {

	s.file.Seek(0, io.SeekStart)
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

func (s *store) Get(key string) []storage.URLRecord {
	var result []storage.URLRecord

	s.file.Seek(0, io.SeekStart)
	for {
		b, _, err := s.reader.ReadLine()
		if err != nil {
			break
		}
		record := strings.Split(string(b), ",")
		if record[0] == key {
			result = append(result, storage.URLRecord{
				Short:    key,
				Original: record[1],
			})
		}
	}
	return result
}

func (s *store) Save(short, original, userID string) {
	s.file.Seek(0, io.SeekEnd)
	record := []string{short, original}

	_, err := s.writer.WriteString(strings.Join(record, ",") + "\n")
	if err != nil {
		log.Println(err)
	}
	s.writer.Flush()
}

func (s *store) IsExist(key string) bool {
	_, err := s.First(key)
	return err == nil
}

func (s *store) Lock() {
	s.fileMutex.Lock()
}
func (s *store) Unlock() {
	s.fileMutex.Unlock()
}
