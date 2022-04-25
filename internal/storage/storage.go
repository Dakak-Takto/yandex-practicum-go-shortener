package storage

type Storage interface {
	GetByShort(key string) (URLRecord, error)
	Save(short, original, userID string)
	IsExist(key string) bool
	GetByUID(uid string) ([]URLRecord, error)
	Lock()
	Unlock()
	Ping() error
}

type URLRecord struct {
	Short    string `json:"short_url" db:"short"`
	Original string `json:"original_url" db:"original"`
	UserID   string `json:"-" db:"user_id"`
}
