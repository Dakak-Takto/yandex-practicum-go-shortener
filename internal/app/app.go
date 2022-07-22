// Package app реализует бизнес-логику приложения и содержит HTTP хендлеры
package app

import (
	"fmt"
	"net/url"

	"github.com/sirupsen/logrus"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/random"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/storage"
)

type Application interface {
	MakeShort(original string, userID string) (storage.URLRecord, error)
	GetByOriginal(original string) (storage.URLRecord, error)
	GetByShort(short string) (storage.URLRecord, error)
	SelectByUID(userID string) ([]storage.URLRecord, error)
	Delete(userID string, shorts ...string)
	PingDatabase() error
}

type application struct {
	store   storage.Storage // хранилище
	baseURL string          // базовый URL для сокращенных ссылок
	addr    string          // адрес для HTTP сервера
	log     *logrus.Logger  // логирование
}

// New создает экземпляр Application и применяет опции
func New(opts ...Option) Application {
	app := application{}
	for _, o := range opts {
		o(&app)
	}
	return &app
}

//Application option declaration

type Option func(app *application)

//WithStorage add storage to application
func WithStorage(storage storage.Storage) Option {
	return func(app *application) {
		app.store = storage
	}
}

//WithLogger set logger
func WithLogger(log *logrus.Logger) Option {
	return func(app *application) {
		app.log = log
	}
}

//WithBaseURL change application base_url
func WithBaseURL(baseURL string) Option {
	return func(app *application) {
		app.baseURL = baseURL
	}
}

//WithAddr change http server addr
func WithAddr(addr string) Option {
	return func(app *application) {
		app.addr = addr
	}
}

// возвращает сокращенную ссылку.
// original - строка, содержащая оригинальный URL
// userID - строка, содержащая идентификатор пользователя
func (app *application) MakeShort(original string, userID string) (storage.URLRecord, error) {
	parsedURL, err := url.ParseRequestURI(original)
	if err != nil {
		return storage.URLRecord{}, fmt.Errorf("no valid url found")
	}

	short := random.String(8)
	if err := app.store.Save(short, parsedURL.String(), userID); err != nil {
		return storage.URLRecord{}, err
	}

	return storage.URLRecord{
		Original: original,
		Short:    short,
		UserID:   userID,
	}, nil
}

func (app *application) GetByOriginal(original string) (storage.URLRecord, error) {
	if url, err := app.store.GetByOriginal(original); err != nil {
		return storage.URLRecord{}, err
	} else {
		return url, nil
	}
}

func (app *application) GetByShort(short string) (storage.URLRecord, error) {
	if url, err := app.store.GetByShort(short); err != nil {
		return storage.URLRecord{}, err
	} else {
		return url, nil
	}
}

func (app *application) SelectByUID(userID string) ([]storage.URLRecord, error) {
	if urls, err := app.store.SelectByUID(userID); err != nil {
		return nil, err
	} else {
		return urls, nil
	}
}

func (app *application) PingDatabase() error {
	if err := app.store.Ping(); err != nil {
		return err
	} else {
		return nil
	}
}

func (app *application) Delete(userID string, shorts ...string) {
	app.store.Delete(userID, shorts...)
}
