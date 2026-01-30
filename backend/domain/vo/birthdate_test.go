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
		{"valid birthdate 1990-01-01", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), nil},
		{"valid birthdate 2000-12-31", time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC), nil},
		// 異常系
		{"future date", tomorrow, domainErrors.ErrBirthDateMustBePast},
		{"today", today, domainErrors.ErrBirthDateMustBePast},
		{"too old 151 years ago", yearsAgo151, domainErrors.ErrBirthDateTooOld},
		// 境界値
		{"yesterday", yesterday, nil},
		{"exactly 150 years ago", yearsAgo150, nil},
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
		{"born 1990-01-01, age 34", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), 34},
		{"born 2000-06-15, age 24 (birthday today)", time.Date(2000, 6, 15, 0, 0, 0, 0, time.UTC), 24},
		{"born 2000-06-16, age 23 (birthday tomorrow)", time.Date(2000, 6, 16, 0, 0, 0, 0, time.UTC), 23},
		{"born 2000-06-14, age 24 (birthday yesterday)", time.Date(2000, 6, 14, 0, 0, 0, 0, time.UTC), 24},
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
