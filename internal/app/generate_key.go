package app

import (
	"yandex-practicum-go-shortener/internal/random"
)

func (app *application) generateKey(startLenght int) string {
	var n = startLenght

	for {
		key := random.String(n)
		if app.store.IsExist(key) {
			n = n + 1
			continue
		} else {
			return key
		}
	}

}
