package storage

type Storage interface {
	First(key string) (URLRecord, error)
	Get(key string) []URLRecord
	Save(short, original, userID string)
	IsExist(key string) bool
	GetByUID(uid string) []URLRecord
	Lock()
	Unlock()
	Ping() error
}

type URLRecord struct {
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
	UserID   string `json:"-"`
}
