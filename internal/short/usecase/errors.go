package usecase

import "errors"

var (
	ErrDuplicate = errors.New("error duplicate")
	ErrNotFound  = errors.New("error not found")
)
