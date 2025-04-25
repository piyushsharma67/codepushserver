package v2

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccessDenied = errors.New("access denied")
) 