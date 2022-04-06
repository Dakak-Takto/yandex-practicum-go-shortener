package inmem

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"sync"
	"yandex-practicum-go-shortener/internal/storage"
)

type store struct {
	urlsMutex sync.Mutex
	urls      map[string]string
	dumpFile  *os.File
	encoder   *json.Encoder
	decoder   *json.Decoder
}

var _ storage.Storage = (*store)(nil)

func New(opts ...option) storage.Storage {
	var s = store{
		urls: make(map[string]string),
	}
	for _, o := range opts {
		o(&s)
	}
	return &s
}

func (s *store) Get(key string) (value string, err error) {
	if value, ok := s.urls[key]; ok {
		return value, nil
	}
	return "", errors.New("not found")
}

func (s *store) Set(key, value string) error {
	s.urls[key] = value
	if s.dumpFile != nil {
		if err := s.dumpToFile(); err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (s *store) IsExist(key string) (isExists bool) {
	_, isExists = s.urls[key]
	return isExists
}

func (s *store) Lock() {
	s.urlsMutex.Lock()
}

func (s *store) Unlock() {
	s.urlsMutex.Unlock()
}

func (s *store) dumpToFile() error {
	s.dumpFile.Seek(0, io.SeekStart)
	return s.encoder.Encode(s.urls)
}

func (s *store) loadFromFile() error {
	s.dumpFile.Seek(0, io.SeekStart)
	return s.decoder.Decode(&s.urls)
}

type option func(s *store)

func WithDumpFile(filepath string) option {
	return func(s *store) {
		file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR|os.O_SYNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
		s.dumpFile = file
		s.encoder = json.NewEncoder(file)
		s.decoder = json.NewDecoder(file)
		s.loadFromFile()
	}
}
