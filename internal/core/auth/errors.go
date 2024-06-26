package auth

import "errors"

var (
	ErrInvalidUserID     = errors.New("invalid user id")
	ErrInvalidUser       = errors.New("invalid email or password")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidClaims     = errors.New("invalid claims")
)
