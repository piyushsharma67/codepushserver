package utils

import "errors"

var (
	ErrAccessDenied = errors.New("access denied")
	ErrNotFound     = errors.New("not found")
	ErrInternal     = errors.New("internal server error")
)
