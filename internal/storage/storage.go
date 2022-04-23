package storage

type Storage interface {
	First(key string) (URLRecord, error)
	Get(key string) []URLRecord
	Save(short, original, userID string)
	IsExist(key string) bool
	Lock()
	Unlock()
}

type URLRecord struct {
	Short, Original, UserID string
}
