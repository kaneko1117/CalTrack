package dto

import (
	"caltrack/domain/entity"
)

type RegisterUserResponse struct {
	UserID string `json:"userId" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// UpdateProfileResponse はプロフィール更新レスポンスDTO
type UpdateProfileResponse struct {
	UserID        string  `json:"userId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Nickname      string  `json:"nickname" example:"NewNickname"`
	Height        float64 `json:"height" example:"175.0"`
	Weight        float64 `json:"weight" example:"70.5"`
	ActivityLevel string  `json:"activityLevel" example:"moderate"`
}

// NewUpdateProfileResponse はEntityからレスポンスDTOを生成する
func NewUpdateProfileResponse(user *entity.User) UpdateProfileResponse {
	return UpdateProfileResponse{
		UserID:        user.ID().String(),
		Nickname:      user.Nickname().String(),
		Height:        user.Height().Cm(),
		Weight:        user.Weight().Kg(),
		ActivityLevel: user.ActivityLevel().String(),
	}
}

// GetProfileResponse はプロフィール取得レスポンスDTO
type GetProfileResponse struct {
	Email         string  `json:"email" example:"user@example.com"`
	Nickname      string  `json:"nickname" example:"John"`
	Weight        float64 `json:"weight" example:"70.5"`
	Height        float64 `json:"height" example:"175.0"`
	BirthDate     string  `json:"birthDate" example:"1990-01-15"`
	Gender        string  `json:"gender" example:"male"`
	ActivityLevel string  `json:"activityLevel" example:"moderate"`
}

// NewGetProfileResponse はEntityからレスポンスDTOを生成する
func NewGetProfileResponse(user *entity.User) GetProfileResponse {
	return GetProfileResponse{
		Email:         user.Email().String(),
		Nickname:      user.Nickname().String(),
		Weight:        user.Weight().Kg(),
		Height:        user.Height().Cm(),
		BirthDate:     user.BirthDate().Time().Format("2006-01-02"),
		Gender:        user.Gender().String(),
		ActivityLevel: user.ActivityLevel().String(),
	}
}
