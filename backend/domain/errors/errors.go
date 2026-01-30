package errors

import "errors"

var (
	ErrEmailRequired       = errors.New("email is required")
	ErrInvalidEmailFormat  = errors.New("invalid email format")
	ErrEmailTooLong        = errors.New("email must be 254 characters or less")
)
