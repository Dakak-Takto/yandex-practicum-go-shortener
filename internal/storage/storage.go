package storage

type Storage interface {
	First(key string) (Entity, error)
	Get(key string) []Entity
	Insert(key, value string)
	IsExist(key string) bool
	Lock()
	Unlock()
}

type Entity struct {
	Key, Value string
}
