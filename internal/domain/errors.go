package domain

import "errors"

var (
	ErrKeyNotFound     = errors.New("key not found")
	ErrKeyIsNotValid   = errors.New("key is not valid")
	ErrValueIsNotValid = errors.New("value is not valid")
)
