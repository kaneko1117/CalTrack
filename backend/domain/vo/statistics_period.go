package vo

import (
	domainErrors "caltrack/domain/errors"
)

const (
	StatisticsPeriodWeek  = "week"
	StatisticsPeriodMonth = "month"
)

var validStatisticsPeriods = map[string]bool{
	StatisticsPeriodWeek:  true,
	StatisticsPeriodMonth: true,
}

// StatisticsPeriod は統計期間を表すValue Object
type StatisticsPeriod struct {
	value string
}

// NewStatisticsPeriod は新しいStatisticsPeriodを生成する
// 空文字の場合はデフォルトでweekを設定する
// week または month のみ許可する
func NewStatisticsPeriod(value string) (StatisticsPeriod, error) {
	if value == "" {
		return StatisticsPeriod{value: StatisticsPeriodWeek}, nil
	}
	if !validStatisticsPeriods[value] {
		return StatisticsPeriod{}, domainErrors.ErrInvalidStatisticsPeriod
	}
	return StatisticsPeriod{value: value}, nil
}

// String は統計期間の文字列表現を返す
func (p StatisticsPeriod) String() string {
	return p.value
}

// Days は統計期間の日数を返す
func (p StatisticsPeriod) Days() int {
	if p.value == StatisticsPeriodWeek {
		return 7
	}
	return 30
}

// IsWeek は期間がweekかどうかを返す
func (p StatisticsPeriod) IsWeek() bool {
	return p.value == StatisticsPeriodWeek
}

// IsMonth は期間がmonthかどうかを返す
func (p StatisticsPeriod) IsMonth() bool {
	return p.value == StatisticsPeriodMonth
}
