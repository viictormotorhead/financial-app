package auth

import "errors"

var ErrUnauthorized = errors.New("unauthorized")

var ErrInvalidCredentials = errors.New("invalid username or password")
