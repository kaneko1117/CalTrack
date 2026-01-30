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
		{"sedentary", "sedentary", "sedentary", nil},
		{"light", "light", "light", nil},
		{"moderate", "moderate", "moderate", nil},
		{"active", "active", "active", nil},
		{"veryActive", "veryActive", "veryActive", nil},
		// 異常系
		{"empty", "", "", domainErrors.ErrInvalidActivityLevel},
		{"invalid", "unknown", "", domainErrors.ErrInvalidActivityLevel},
		{"wrong case", "Sedentary", "", domainErrors.ErrInvalidActivityLevel},
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
		{"sedentary", "sedentary", 1.2},
		{"light", "light", 1.375},
		{"moderate", "moderate", 1.55},
		{"active", "active", 1.725},
		{"veryActive", "veryActive", 1.9},
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
