package storage

import (
	"errors"
	"sync"
	"time"
)

type Storage interface {
	Get(key string) (value string, err error)
	Set(key, value string)
	IsExist(key string) bool
	Lock()
	Unlock()
}

type storage struct {
	urlsMutex sync.Mutex
	urls      map[string]string
}

var _ Storage = (*storage)(nil)

func New() Storage {
	return &storage{
		urls: make(map[string]string),
	}
}

func (s *storage) Get(key string) (value string, err error) {
	if value, ok := s.urls[key]; ok {
		return value, nil
	}
	return "", errors.New("not found")
}

func (s *storage) Set(key, value string) {
	s.urls[key] = value
}

func (s *storage) IsExist(key string) (isExists bool) {
	_, isExists = s.urls[key]
	time.Sleep(time.Second * 2)
	return isExists
}

func (s *storage) Lock() {
	s.urlsMutex.Lock()
}

func (s *storage) Unlock() {
	s.urlsMutex.Unlock()
}
