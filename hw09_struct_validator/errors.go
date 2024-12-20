package hw09structvalidator

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotStruct     = errors.New("input is not a struct")
	ErrInvalidTag    = errors.New("invalid validation tag")
	ErrInvalidLength = errors.New("invalid length")
	ErrInvalidMin    = errors.New("value is less than minimum")
	ErrInvalidMax    = errors.New("value exceeds maximum")
	ErrInvalidRegexp = errors.New("value does not match regular expression")
	ErrInvalidIn     = errors.New("value is not in the allowed set")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, err := range v {
		sb.WriteString(fmt.Sprintf("Field %s: %s\n", err.Field, err.Err))
	}
	return sb.String()
}
