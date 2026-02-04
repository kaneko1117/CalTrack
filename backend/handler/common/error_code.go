package common

const (
	CodeInvalidRequest     = "INVALID_REQUEST"
	CodeValidationError    = "VALIDATION_ERROR"
	CodeEmailAlreadyExists = "EMAIL_ALREADY_EXISTS"
	CodeInternalError      = "INTERNAL_ERROR"
	CodeNotFound           = "NOT_FOUND"

	// 認証関連エラーコード
	CodeInvalidCredentials = "INVALID_CREDENTIALS"
	CodeUnauthorized       = "UNAUTHORIZED"
	CodeSessionExpired     = "SESSION_EXPIRED"

	// 画像解析関連エラーコード
	CodeNoFoodDetected      = "NO_FOOD_DETECTED"
	CodeImageAnalysisFailed = "IMAGE_ANALYSIS_FAILED"
)
