package storage

import "net/url"

type Repository interface {
	//Search original URL by short key. Return URL if exist or error
	Get(key string) (url.URL, error)

	//Create new record and return short key
	Create(location url.URL) (string, error)

	//Check is short key in repository
	IsExist(key string) bool
}
