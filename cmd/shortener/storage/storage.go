package storage

import (
	"errors"
	"fmt"
	"hash/crc32"
	"math/rand"
	"sync"
	"time"
)

var storage = make(map[string]string)
var mu sync.Mutex
var saveAttemps = 10

func Save(value string) (string, error) {
	for i := 0; i < saveAttemps; i++ {
		key := generateKey(value)
		if !isExist(key) {
			mu.Lock()
			storage[key] = value
			mu.Unlock()
			return key, nil
		}
	}
	return "", errors.New("free key not found")
}

func Get(key string) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	if value, ok := storage[key]; ok {
		return value, nil
	} else {
		return "", errors.New("item not found")
	}
}

func isExist(key string) bool {
	_, ok := storage[key]
	return ok
}

func generateKey(str string) string {
	b := []byte(str)
	hash := crc32.ChecksumIEEE(b)
	hash += uint32(time.Now().UnixMicro()) + rand.Uint32()
	return fmt.Sprintf("%x", hash)
}

func Len() int {
	return len(storage)
}
