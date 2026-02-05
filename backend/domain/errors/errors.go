package errors

import "errors"

var (
	ErrEmailRequired          = errors.New("email is required")
	ErrInvalidEmailFormat     = errors.New("invalid email format")
	ErrEmailTooLong           = errors.New("email must be 254 characters or less")
	ErrPasswordRequired       = errors.New("password is required")
	ErrPasswordTooShort       = errors.New("password must be at least 8 characters")
	ErrNicknameRequired       = errors.New("nickname is required")
	ErrNicknameTooLong        = errors.New("nickname must be 50 characters or less")
	ErrWeightMustBePositive   = errors.New("weight must be positive")
	ErrWeightTooHeavy         = errors.New("weight must be 500kg or less")
	ErrHeightMustBePositive   = errors.New("height must be positive")
	ErrHeightTooTall          = errors.New("height must be 300cm or less")
	ErrBirthDateMustBePast    = errors.New("birth date must be in the past")
	ErrBirthDateTooOld        = errors.New("birth date must be within 150 years")
	ErrEatenAtMustNotBeFuture = errors.New("eaten at must not be in the future")
	ErrCaloriesMustBePositive = errors.New("calories must be positive")
	ErrInvalidGender          = errors.New("gender must be male, female, or other")
	ErrInvalidActivityLevel   = errors.New("activity level must be sedentary, light, moderate, active, or veryActive")

	// Usecase errors
	ErrEmailAlreadyExists = errors.New("email already exists")

	// Auth errors
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrSessionNotFound    = errors.New("session not found")
	ErrUserNotFound       = errors.New("user not found")

	// Session errors
	ErrSessionIDGenerationFailed = errors.New("failed to generate session id")
	ErrInvalidSessionID          = errors.New("invalid session id")
	ErrSessionExpired            = errors.New("session has expired")

	// UUID errors
	ErrUUIDRequired      = errors.New("uuid is required")
	ErrInvalidUUIDFormat = errors.New("invalid uuid format")

	// ID errors
	ErrInvalidUserID        = errors.New("invalid user id")
	ErrInvalidRecordID      = errors.New("invalid record id")
	ErrInvalidRecordItemID  = errors.New("invalid record item id")
	ErrInvalidRecordPfcID   = errors.New("invalid record pfc id")
	ErrInvalidAdviceCacheID = errors.New("invalid advice cache id")

	// Record Item errors
	ErrItemNameRequired = errors.New("item name is required")

	// Statistics errors
	ErrInvalidStatisticsPeriod = errors.New("statistics period must be week or month")

	// 画像解析関連エラー
	ErrImageDataRequired   = errors.New("画像データは必須です")
	ErrMimeTypeRequired    = errors.New("MIMEタイプは必須です")
	ErrNoFoodDetected      = errors.New("食品を検出できませんでした")
	ErrImageAnalysisFailed = errors.New("画像解析に失敗しました")
)
