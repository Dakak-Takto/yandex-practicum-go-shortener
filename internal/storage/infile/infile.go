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

func (s *store) First(key string) (storage.Entity, error) {

	s.file.Seek(0, io.SeekStart)
	for {
		b, _, err := s.reader.ReadLine()
		if err != nil {
			break
		}
		record := strings.Split(string(b), ",")
		if key == record[0] {
			return storage.Entity{
				Key:   record[0],
				Value: record[1],
			}, nil
		}
	}
	return storage.Entity{}, errors.New("errNotFound")
}

func (s *store) Get(key string) []storage.Entity {
	var result []storage.Entity

	s.file.Seek(0, io.SeekStart)
	for {
		b, _, err := s.reader.ReadLine()
		if err != nil {
			break
		}
		record := strings.Split(string(b), ",")
		if record[0] == key {
			result = append(result, storage.Entity{
				Key:   key,
				Value: record[1],
			})
		}
	}
	return result
}

func (s *store) Insert(key, value string) {
	s.file.Seek(0, io.SeekEnd)
	record := []string{key, value}

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
