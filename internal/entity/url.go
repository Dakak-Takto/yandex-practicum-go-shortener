package entity

import (
	"context"
	"errors"
)

type (
	Short string

	URL struct {
		Short    string `json:"short_url" db:"short"`
		Original string `json:"original_url" db:"original"`
		UserID   string `json:"-" db:"user_id"`
		Deleted  bool   `json:"-" db:"deleted"`
	}

	URLRepository interface {
		GetByShort(ctx context.Context, short string) (*URL, error)
		GetByOriginal(ctx context.Context, original string) (*URL, error)
		Save(ctx context.Context, url *URL) error
		SelectByUserID(ctx context.Context, uid string) ([]*URL, error)
		Delete(uid string, keys ...string)
	}

	URLUsecase interface {
		Create(original string, userID string) (*URL, error)
		GetByShort(short string) (*URL, error)
		GetByOriginal(original string) (*URL, error)
		UserURLs(userID string) ([]*URL, error)
		Delete(userID string, shorts ...string)
	}
)

var (
	ErrDuplicate = errors.New("duplicate")
	ErrNotFound  = errors.New("not found")
)
