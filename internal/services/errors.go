package services

import "errors"

var (
	ErrURLAlreadyExists = errors.New("alias already exists")
	ErrInvalidInput     = errors.New("invalid input")
	ErrURLNotFound      = errors.New("invalid input")
)
