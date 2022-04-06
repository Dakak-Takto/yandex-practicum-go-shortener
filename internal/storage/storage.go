package storage

type Storage interface {
	Get(key string) (value string, err error)
	Set(key, value string) error
	IsExist(key string) bool
	Lock()
	Unlock()
}
