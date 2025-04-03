package errors

import "errors"

var (
	ErrURLExists = errors.New("URL уже существует")
	ErrDeleted   = errors.New("удалено")
	ErrNotFound  = errors.New("URL не найден")
)
