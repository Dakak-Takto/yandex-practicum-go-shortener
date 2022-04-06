package app

import (
	"errors"
	"time"
	"yandex-practicum-go-shortener/internal/random"
)

func (app *application) generateKey(startLenght int) (string, error) {

	var ch = make(chan string)

	go func() {
		var n = startLenght
		key := random.String(n)
		if app.store.IsExist(key) {
			n = n + 1
		} else {
			ch <- key
			return
		}
	}()

	for {
		select {
		case <-time.Tick(time.Second):
			return "", errors.New("timeout")
		case key := <-ch:
			return key, nil
		}
	}

}
