package errors

import "errors"

// Error variables used throughout the application to represent error
var (
	// ErrURLExists indicates that the given URL already exists
	ErrURLExists = errors.New("URL уже существует")
	// ErrDeleted indicates that the resource has been deleted
	ErrDeleted = errors.New("удалено")
	// ErrNotFound indicates that the specified URL was not found
	ErrNotFound = errors.New("URL не найден")
)
