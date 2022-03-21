package storage

import "errors"

type Links map[string]string

func CreateNew() *Links {
	links := make(Links)
	return &links
}

func (l Links) Get(id string) (string, error) {
	if link, ok := l[id]; ok {
		return link, nil
	} else {
		err := errors.New("not found")
		return "", err
	}
}

func (l Links) Set(short string, long string) {
	l[short] = long
}
