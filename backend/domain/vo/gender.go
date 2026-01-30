package vo

import (
	domainErrors "caltrack/domain/errors"
)

const (
	GenderMale   = "male"
	GenderFemale = "female"
	GenderOther  = "other"
)

var validGenders = map[string]bool{
	GenderMale:   true,
	GenderFemale: true,
	GenderOther:  true,
}

type Gender struct {
	value string
}

func NewGender(value string) (Gender, error) {
	if !validGenders[value] {
		return Gender{}, domainErrors.ErrInvalidGender
	}
	return Gender{value: value}, nil
}

func (g Gender) String() string {
	return g.value
}
