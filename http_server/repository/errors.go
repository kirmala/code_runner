package repository

import "errors"

var (
	NotFound      = errors.New("key not found")
	AlreadyExists = errors.New("key already exists")
)
