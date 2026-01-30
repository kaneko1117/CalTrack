package vo

import (
	domainErrors "caltrack/domain/errors"
)

const maxHeight = 300.0

type Height struct {
	value float64
}

func NewHeight(cm float64) (Height, error) {
	if cm <= 0 {
		return Height{}, domainErrors.ErrHeightMustBePositive
	}
	if cm > maxHeight {
		return Height{}, domainErrors.ErrHeightTooTall
	}
	return Height{value: cm}, nil
}

func (h Height) Cm() float64 {
	return h.value
}

func (h Height) Meters() float64 {
	return h.value / 100
}
