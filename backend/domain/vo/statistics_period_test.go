package vo_test

import (
	"testing"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
)

func TestNewStatisticsPeriod(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantPeriod string
		wantErr    error
	}{
		// 正常系
		{"weekは有効", "week", "week", nil},
		{"monthは有効", "month", "month", nil},
		{"空文字はデフォルトでweek", "", "week", nil},
		// 異常系
		{"無効な値はエラー", "year", "", domainErrors.ErrInvalidStatisticsPeriod},
		{"大文字始まりはエラー", "Week", "", domainErrors.ErrInvalidStatisticsPeriod},
		{"dayはエラー", "day", "", domainErrors.ErrInvalidStatisticsPeriod},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vo.NewStatisticsPeriod(tt.input)

			if err != tt.wantErr {
				t.Errorf("NewStatisticsPeriod(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if err == nil && got.String() != tt.wantPeriod {
				t.Errorf("NewStatisticsPeriod(%q).String() = %v, want %v", tt.input, got.String(), tt.wantPeriod)
			}
		})
	}
}

func TestStatisticsPeriod_Days(t *testing.T) {
	tests := []struct {
		name     string
		period   string
		wantDays int
	}{
		{"weekの日数は7", "week", 7},
		{"monthの日数は30", "month", 30},
		{"空文字のデフォルト（week）の日数は7", "", 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp, _ := vo.NewStatisticsPeriod(tt.period)
			if got := sp.Days(); got != tt.wantDays {
				t.Errorf("StatisticsPeriod.Days() = %v, want %v", got, tt.wantDays)
			}
		})
	}
}

func TestStatisticsPeriod_IsWeek(t *testing.T) {
	tests := []struct {
		name       string
		period     string
		wantIsWeek bool
	}{
		{"weekはtrue", "week", true},
		{"monthはfalse", "month", false},
		{"空文字のデフォルト（week）はtrue", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp, _ := vo.NewStatisticsPeriod(tt.period)
			if got := sp.IsWeek(); got != tt.wantIsWeek {
				t.Errorf("StatisticsPeriod.IsWeek() = %v, want %v", got, tt.wantIsWeek)
			}
		})
	}
}

func TestStatisticsPeriod_IsMonth(t *testing.T) {
	tests := []struct {
		name        string
		period      string
		wantIsMonth bool
	}{
		{"weekはfalse", "week", false},
		{"monthはtrue", "month", true},
		{"空文字のデフォルト（week）はfalse", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp, _ := vo.NewStatisticsPeriod(tt.period)
			if got := sp.IsMonth(); got != tt.wantIsMonth {
				t.Errorf("StatisticsPeriod.IsMonth() = %v, want %v", got, tt.wantIsMonth)
			}
		})
	}
}
