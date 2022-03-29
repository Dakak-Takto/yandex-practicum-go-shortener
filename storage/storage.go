package storage

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
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
	for {
		t := uint32(time.Now().UnixMicro())
		r := rand.Uint32()
		key := fmt.Sprintf("%x", t+r)
		if keyIsExist(key) {
			continue
		}
		return key
	}
}
