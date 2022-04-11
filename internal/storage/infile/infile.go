//Used for file storage. Making file and write url per line
package infile

import (
	"bufio"
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

func (s *store) Get(key string) (value string, err error) {
	s.file.Seek(0, io.SeekStart)
	for {
		b, _, err := s.reader.ReadLine()
		if err != nil {
			return "", err
		}
		record := strings.Split(string(b), ",")
		if record[0] == key {
			return record[1], nil
		}
	}
}

func (s *store) Set(key, value string) (err error) {
	s.file.Seek(0, io.SeekEnd)
	record := []string{key, value}

	_, err = s.writer.WriteString(strings.Join(record, ",") + "\n")
	if err != nil {
		return err
	}
	return s.writer.Flush()
}

func (s *store) IsExist(key string) bool {
	_, err := s.Get(key)
	return err == nil
}

func (s *store) Lock() {
	s.fileMutex.Lock()
}
func (s *store) Unlock() {
	s.fileMutex.Unlock()
}
