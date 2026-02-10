package dto

import "caltrack/domain/vo"

// LoginRequest はログインリクエストDTO
type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

// ToDomain はリクエストをドメインのVOに変換する
func (r LoginRequest) ToDomain() (vo.Email, vo.Password, error) {
	email, err := vo.NewEmail(r.Email)
	if err != nil {
		return vo.Email{}, vo.Password{}, err
	}

	password, err := vo.NewPassword(r.Password)
	if err != nil {
		return vo.Email{}, vo.Password{}, err
	}

	return email, password, nil
}
