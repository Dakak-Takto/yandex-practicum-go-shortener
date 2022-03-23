package storage

import (
	"errors"
	"sync"
)

var storage = make(map[string]string)
var mu sync.Mutex

func Get(key string) (string, error) {
	if v, ok := storage[key]; ok {
		return v, nil
	}

	return "", errors.New("not found")
}

func Set(key, value string) error {
	if !isExist(key) {
		storage[key] = value
		return nil
	}

	return errors.New("key already exists")
}

func isExist(key string) bool {
	_, isExists := storage[key]

	return isExists
}

func Lock() {
	mu.Lock()
}

func UnLock() {
	mu.Unlock()
}
