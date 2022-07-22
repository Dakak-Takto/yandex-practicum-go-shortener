//Package storage модуль хранения записей сокращенных URL
package storage

import "errors"

//go:generate mockgen -destination=mocks/storage.go -package=mocks github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/storage Storage
type Storage interface {
	//GetByShort search and return URLRecord by short key
	GetByShort(key string) (URLRecord, error)
	//GetByOriginal search and return []URLRecord by original URL
	GetByOriginal(original string) (URLRecord, error)
	//Save store new URLRecord
	Save(short, original, userID string) error
	//SelectByUID return []URLRecord by userID
	SelectByUID(uid string) ([]URLRecord, error)
	//Ping check database connectivity
	Ping() error
	//Delete remove urls by key and userID
	Delete(uid string, keys ...string)
}

type URLRecord struct {
	Short    string `json:"short_url" db:"short"`       //short url key
	Original string `json:"original_url" db:"original"` //original URL
	UserID   string `json:"-" db:"user_id"`             //User id
	Deleted  bool   `json:"-" db:"deleted"`             //deleted flag
}

var (
	ErrDuplicate = errors.New("error duplicate") // original url exists in storage
)
