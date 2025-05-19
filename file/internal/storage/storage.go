package storage

import (
	"errors"
	"io"
)

// Storage defines the behavior for file operations
// Implementations may be of the time local disk, or cloud storage, etc
type Storage interface {
	Save(path string, file io.Reader) error
}

var ErrDirPermDenied = errors.New("cannot create directory, permission denied")
