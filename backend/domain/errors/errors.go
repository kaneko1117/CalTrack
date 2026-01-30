package errors

import "errors"

var (
	ErrEmailRequired       = errors.New("email is required")
	ErrInvalidEmailFormat  = errors.New("invalid email format")
	ErrEmailTooLong        = errors.New("email must be 254 characters or less")
	ErrInvalidUserID       = errors.New("invalid user id: must be a valid UUID")
	ErrPasswordRequired    = errors.New("password is required")
	ErrPasswordTooShort    = errors.New("password must be at least 8 characters")
	ErrNicknameRequired    = errors.New("nickname is required")
	ErrNicknameTooLong     = errors.New("nickname must be 50 characters or less")
	ErrWeightMustBePositive = errors.New("weight must be positive")
	ErrWeightTooHeavy       = errors.New("weight must be 500kg or less")
	ErrHeightMustBePositive = errors.New("height must be positive")
	ErrHeightTooTall        = errors.New("height must be 300cm or less")
	ErrBirthDateMustBePast  = errors.New("birth date must be in the past")
	ErrBirthDateTooOld      = errors.New("birth date must be within 150 years")
	ErrInvalidGender        = errors.New("gender must be male, female, or other")
	ErrInvalidActivityLevel = errors.New("activity level must be sedentary, light, moderate, active, or veryActive")

	// Usecase errors
	ErrEmailAlreadyExists = errors.New("email already exists")
)
