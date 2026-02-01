package dto

import "caltrack/usecase"

// LoginResponse はログインレスポンスDTO
type LoginResponse struct {
	UserID   string `json:"userId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email    string `json:"email" example:"user@example.com"`
	Nickname string `json:"nickname" example:"John"`
}

// NewLoginResponse はUsecaseの出力からレスポンスDTOを生成する
func NewLoginResponse(output *usecase.LoginOutput) LoginResponse {
	return LoginResponse{
		UserID:   output.User.ID().String(),
		Email:    output.User.Email().String(),
		Nickname: output.User.Nickname().String(),
	}
}
