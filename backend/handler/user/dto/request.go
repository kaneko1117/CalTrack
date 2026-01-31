package dto

import (
	"time"

	"caltrack/domain/entity"
)

type RegisterUserRequest struct {
	Email         string  `json:"email"`
	Password      string  `json:"password"`
	Nickname      string  `json:"nickname"`
	Weight        float64 `json:"weight"`
	Height        float64 `json:"height"`
	BirthDate     string  `json:"birthDate"` // "2006-01-02"形式
	Gender        string  `json:"gender"`
	ActivityLevel string  `json:"activityLevel"`
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
