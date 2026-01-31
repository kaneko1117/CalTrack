package dto

import "caltrack/usecase"

// LoginResponse はログインレスポンスDTO
type LoginResponse struct {
	UserID   string `json:"userId"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

// NewLoginResponse はUsecaseの出力からレスポンスDTOを生成する
func NewLoginResponse(output *usecase.LoginOutput) LoginResponse {
	return LoginResponse{
		UserID:   output.User.ID().String(),
		Email:    output.User.Email().String(),
		Nickname: output.User.Nickname().String(),
	}
}
