package vo

import (
	domainErrors "caltrack/domain/errors"
)

const maxWeight = 500.0

type Weight struct {
	value float64
}

func NewWeight(kg float64) (Weight, error) {
	if kg <= 0 {
		return Weight{}, domainErrors.ErrWeightMustBePositive
	}
	if kg > maxWeight {
		return Weight{}, domainErrors.ErrWeightTooHeavy
	}
	return Weight{value: kg}, nil
}

func (w Weight) Kg() float64 {
	return w.value
}
