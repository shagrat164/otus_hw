package storage

import "errors"

var (
	ErrEventAlreadyExist = errors.New("event already exists")
	ErrEventNotFound     = errors.New("event not found")
)
