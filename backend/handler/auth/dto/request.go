package dto

import "caltrack/usecase"

// LoginRequest はログインリクエストDTO
type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

// ToInput はリクエストをUsecase入力に変換する
func (r LoginRequest) ToInput() usecase.LoginInput {
	return usecase.LoginInput{
		Email:    r.Email,
		Password: r.Password,
	}
}
