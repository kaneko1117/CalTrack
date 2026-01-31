package vo

import (
	"testing"
	"time"

	domainErrors "caltrack/domain/errors"
)

func TestNewBirthDate(t *testing.T) {
	// Fix "now" for consistent testing
	fixedNow := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	nowFunc = func() time.Time { return fixedNow }
	defer func() { nowFunc = time.Now }()

	today := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	yesterday := today.AddDate(0, 0, -1)
	tomorrow := today.AddDate(0, 0, 1)
	yearsAgo150 := today.AddDate(-150, 0, 0)
	yearsAgo151 := today.AddDate(-151, 0, 0)

	tests := []struct {
		name    string
		input   time.Time
		wantErr error
	}{
		// 正常系
		{"正常な生年月日1990-01-01", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), nil},
		{"正常な生年月日2000-12-31", time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC), nil},
		// 異常系
		{"未来の日付はエラー", tomorrow, domainErrors.ErrBirthDateMustBePast},
		{"今日の日付はエラー", today, domainErrors.ErrBirthDateMustBePast},
		{"151年前はエラー", yearsAgo151, domainErrors.ErrBirthDateTooOld},
		// 境界値
		{"昨日は有効", yesterday, nil},
		{"ちょうど150年前は有効", yearsAgo150, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewBirthDate(tt.input)

			if err != tt.wantErr {
				t.Errorf("NewBirthDate(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestBirthDate_Age(t *testing.T) {
	// Fix "now" for consistent testing
	fixedNow := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	nowFunc = func() time.Time { return fixedNow }
	defer func() { nowFunc = time.Now }()

	tests := []struct {
		name      string
		birthDate time.Time
		wantAge   int
	}{
		{"1990-01-01生まれは34歳", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), 34},
		{"2000-06-15生まれは24歳（誕生日当日）", time.Date(2000, 6, 15, 0, 0, 0, 0, time.UTC), 24},
		{"2000-06-16生まれは23歳（誕生日前日）", time.Date(2000, 6, 16, 0, 0, 0, 0, time.UTC), 23},
		{"2000-06-14生まれは24歳（誕生日翌日）", time.Date(2000, 6, 14, 0, 0, 0, 0, time.UTC), 24},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bd, _ := NewBirthDate(tt.birthDate)
			if got := bd.Age(); got != tt.wantAge {
				t.Errorf("BirthDate.Age() = %v, want %v", got, tt.wantAge)
			}
		})
	}
}
