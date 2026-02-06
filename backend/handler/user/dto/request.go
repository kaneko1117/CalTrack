package dto

import (
	"time"

	"caltrack/domain/entity"
	"caltrack/usecase"
)

type RegisterUserRequest struct {
	Email         string  `json:"email" example:"user@example.com"`
	Password      string  `json:"password" example:"password123"`
	Nickname      string  `json:"nickname" example:"John"`
	Weight        float64 `json:"weight" example:"70.5"`
	Height        float64 `json:"height" example:"175.0"`
	BirthDate     string  `json:"birthDate" example:"1990-01-15"` // "2006-01-02"形式
	Gender        string  `json:"gender" example:"male"`
	ActivityLevel string  `json:"activityLevel" example:"moderate"`
}

func (r RegisterUserRequest) ToDomain() (*entity.User, error, []error) {
	birthDate, err := time.Parse("2006-01-02", r.BirthDate)
	if err != nil {
		return nil, err, nil
	}

	user, errs := entity.NewUser(
		r.Email,
		r.Password,
		r.Nickname,
		r.Weight,
		r.Height,
		birthDate,
		r.Gender,
		r.ActivityLevel,
	)

	return user, nil, errs
}

// UpdateProfileRequest はプロフィール更新リクエストDTO
type UpdateProfileRequest struct {
	Nickname      string  `json:"nickname" example:"NewNickname"`
	Height        float64 `json:"height" example:"175.0"`
	Weight        float64 `json:"weight" example:"70.5"`
	ActivityLevel string  `json:"activityLevel" example:"moderate"`
}

// ToInput はUsecaseのUpdateProfileInputに変換する
func (r UpdateProfileRequest) ToInput() usecase.UpdateProfileInput {
	return usecase.UpdateProfileInput{
		Nickname:      r.Nickname,
		Height:        r.Height,
		Weight:        r.Weight,
		ActivityLevel: r.ActivityLevel,
	}
}
