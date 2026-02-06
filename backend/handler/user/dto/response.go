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
