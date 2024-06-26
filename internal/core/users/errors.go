package users

import "errors"

var (
	ErrInvalidUserID     = errors.New("invalid user id")
	ErrInvalidUserEmail  = errors.New("invalid user email")
	ErrUserAlreadyExists = errors.New("user already exists")
)
