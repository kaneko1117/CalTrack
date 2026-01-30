package errors

import "errors"

var (
	ErrEmailRequired       = errors.New("email is required")
	ErrInvalidEmailFormat  = errors.New("invalid email format")
	ErrEmailTooLong        = errors.New("email must be 254 characters or less")
	ErrInvalidUserID       = errors.New("invalid user id: must be a valid UUID")
	ErrPasswordRequired    = errors.New("password is required")
	ErrPasswordTooShort    = errors.New("password must be at least 8 characters")
)
