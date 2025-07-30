package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url now found")
	ErrURLNotExist = errors.New("url exists")
)
