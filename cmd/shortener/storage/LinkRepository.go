package storage

type LinkRepository map[string]string

func (l LinkRepository) CreateNew() LinkRepository {
	return l
}
