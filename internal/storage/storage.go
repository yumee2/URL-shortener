package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url now found")
	ErrURLExist    = errors.New("url exists")
)
