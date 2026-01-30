package vo

import (
	"time"

	domainErrors "caltrack/domain/errors"
)

const maxAgeYears = 150

// nowFunc is used for testing - defaults to time.Now
var nowFunc = time.Now

type BirthDate struct {
	value time.Time
}

func NewBirthDate(date time.Time) (BirthDate, error) {
	now := nowFunc()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	birthDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	if !birthDay.Before(today) {
		return BirthDate{}, domainErrors.ErrBirthDateMustBePast
	}

	minDate := today.AddDate(-maxAgeYears, 0, 0)
	if birthDay.Before(minDate) {
		return BirthDate{}, domainErrors.ErrBirthDateTooOld
	}

	return BirthDate{value: birthDay}, nil
}

func (b BirthDate) Time() time.Time {
	return b.value
}

func (b BirthDate) Age() int {
	now := nowFunc()
	years := now.Year() - b.value.Year()

	if now.YearDay() < b.value.YearDay() {
		years--
	}

	return years
}
