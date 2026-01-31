package vo_test

import (
	"testing"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
)

func TestNewActivityLevel(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantLevel string
		wantErr   error
	}{
		// 正常系
		{"sedentaryは有効", "sedentary", "sedentary", nil},
		{"lightは有効", "light", "light", nil},
		{"moderateは有効", "moderate", "moderate", nil},
		{"activeは有効", "active", "active", nil},
		{"veryActiveは有効", "veryActive", "veryActive", nil},
		// 異常系
		{"空文字はエラー", "", "", domainErrors.ErrInvalidActivityLevel},
		{"無効な値はエラー", "unknown", "", domainErrors.ErrInvalidActivityLevel},
		{"大文字始まりはエラー", "Sedentary", "", domainErrors.ErrInvalidActivityLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vo.NewActivityLevel(tt.input)

			if err != tt.wantErr {
				t.Errorf("NewActivityLevel(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if err == nil && got.String() != tt.wantLevel {
				t.Errorf("NewActivityLevel(%q).String() = %v, want %v", tt.input, got.String(), tt.wantLevel)
			}
		})
	}
}

func TestActivityLevel_Multiplier(t *testing.T) {
	tests := []struct {
		name           string
		level          string
		wantMultiplier float64
	}{
		{"sedentaryの係数は1.2", "sedentary", 1.2},
		{"lightの係数は1.375", "light", 1.375},
		{"moderateの係数は1.55", "moderate", 1.55},
		{"activeの係数は1.725", "active", 1.725},
		{"veryActiveの係数は1.9", "veryActive", 1.9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			al, _ := vo.NewActivityLevel(tt.level)
			if got := al.Multiplier(); got != tt.wantMultiplier {
				t.Errorf("ActivityLevel.Multiplier() = %v, want %v", got, tt.wantMultiplier)
			}
		})
	}
}
