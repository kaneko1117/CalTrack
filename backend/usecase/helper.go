package usecase

import (
	"time"

	"caltrack/domain/helper"
)

func startOfDay(t time.Time) time.Time {
	jstTime := t.In(helper.JST())
	return time.Date(jstTime.Year(), jstTime.Month(), jstTime.Day(), 0, 0, 0, 0, helper.JST())
}

func endOfDay(t time.Time) time.Time {
	return startOfDay(t).AddDate(0, 0, 1)
}
