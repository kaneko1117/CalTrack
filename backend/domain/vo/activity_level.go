package vo

import (
	domainErrors "caltrack/domain/errors"
)

const (
	ActivityLevelSedentary  = "sedentary"
	ActivityLevelLight      = "light"
	ActivityLevelModerate   = "moderate"
	ActivityLevelActive     = "active"
	ActivityLevelVeryActive = "veryActive"
)

var activityMultipliers = map[string]float64{
	ActivityLevelSedentary:  1.2,
	ActivityLevelLight:      1.375,
	ActivityLevelModerate:   1.55,
	ActivityLevelActive:     1.725,
	ActivityLevelVeryActive: 1.9,
}

type ActivityLevel struct {
	value string
}

func NewActivityLevel(value string) (ActivityLevel, error) {
	if _, ok := activityMultipliers[value]; !ok {
		return ActivityLevel{}, domainErrors.ErrInvalidActivityLevel
	}
	return ActivityLevel{value: value}, nil
}

func (a ActivityLevel) String() string {
	return a.value
}

func (a ActivityLevel) Multiplier() float64 {
	return activityMultipliers[a.value]
}
