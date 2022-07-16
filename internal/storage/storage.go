package storage

import "errors"

//go:generate mockgen -destination=mocks/storage.go -package=mocks github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/storage Storage
type Storage interface {
	GetByShort(key string) (URLRecord, error)
	GetByOriginal(original string) (URLRecord, error)
	Save(short, original, userID string) error
	SelectByUID(uid string) ([]URLRecord, error)
	Ping() error
	Delete(uid string, keys ...string)
}

type URLRecord struct {
	Short    string `json:"short_url" db:"short"`
	Original string `json:"original_url" db:"original"`
	UserID   string `json:"-" db:"user_id"`
	Deleted  bool   `json:"-" db:"deleted"`
}

var (
	ErrDuplicate = errors.New("error duplicate")
)
