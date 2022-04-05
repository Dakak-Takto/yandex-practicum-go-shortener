package infile

import (
	"errors"
	"log"
	"net/url"
	"os"
	"sync"
	"syscall"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/random"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/storage"
)

var ErrNotFound = errors.New("not found")

const KeyGeneratorStartLen = 5

type store struct {
	sync.Mutex
	data map[string]url.URL
	file *os.File
}

var _ storage.Repository = (*store)(nil)

//Create and return new infile Repository
func Load(filename string) storage.Repository {
	file, err := os.OpenFile(filename, syscall.O_CREAT|syscall.O_RDWR, 644)
	if err != nil {
		log.Fatal(err)
	}

	var store = store{
		data: make(map[string]url.URL),
		file: file,
	}

	return &store
}

func (s *store) Get(key string) (location url.URL, err error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	if result, ok := s.data[key]; ok {
		location = result
	} else {
		err = ErrNotFound
	}

	return location, err
}

func (s *store) Create(location url.URL) (key string, err error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	for i := 0; ; i++ {
		key = random.String(KeyGeneratorStartLen + i)

		if !s.IsExist(key) {
			break
		}
	}

	s.data[key] = location

	return key, err
}

func (s *store) IsExist(key string) bool {
	_, exist := s.data[key]

	return exist
}
