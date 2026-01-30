package vo

import (
	"regexp"

	domainErrors "caltrack/domain/errors"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type Email struct {
	value string
}

func NewEmail(value string) (Email, error) {
	if value == "" {
		return Email{}, domainErrors.ErrEmailRequired
	}
	if len(value) > 254 {
		return Email{}, domainErrors.ErrEmailTooLong
	}
	if !emailRegex.MatchString(value) {
		return Email{}, domainErrors.ErrInvalidEmailFormat
	}
	return Email{value: value}, nil
}

func (e Email) String() string {
	return e.value
}

func (e Email) Equals(other Email) bool {
	return e.value == other.value
}
