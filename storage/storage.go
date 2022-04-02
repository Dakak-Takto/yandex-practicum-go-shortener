package storage

import (
	"errors"
	"sync"
	"yandex-practicum-go-shortener/internal/random"
)

var storage = make(map[string]string)
var mu sync.Mutex

func GetValueByKey(key string) (value string, err error) {
	if value, ok := storage[key]; ok {
		return value, nil
	}
	return "", errors.New("not found")
}

func SetValueReturnKey(value string) (key string) {
	mu.Lock()
	key = generateUniqueKey()
	storage[key] = value
	mu.Unlock()
	return key
}

func keyIsExist(key string) (isExists bool) {
	_, isExists = storage[key]
	return isExists
}

func generateUniqueKey() (key string) {
	var keyLenght = 5
	for {
		key := random.String(keyLenght)
		if keyIsExist(key) {
			keyLenght = keyLenght + 1
			continue
		}
		return key
	}
}
