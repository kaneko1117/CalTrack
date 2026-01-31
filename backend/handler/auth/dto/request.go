package dto

import "caltrack/usecase"

// LoginRequest はログインリクエストDTO
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ToInput はリクエストをUsecase入力に変換する
func (r LoginRequest) ToInput() usecase.LoginInput {
	return usecase.LoginInput{
		Email:    r.Email,
		Password: r.Password,
	}
}
