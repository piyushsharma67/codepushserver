package errors

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotFound           = errors.New("resource not found")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInternal          = errors.New("internal server error")
) 