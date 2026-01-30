package errors

import "errors"

var (
	ErrInvalidUserID = errors.New("invalid user id: must be a valid UUID")
)
