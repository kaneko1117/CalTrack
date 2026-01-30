package vo_test

import (
	"testing"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
)

func TestNewGender(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantGender string
		wantErr    error
	}{
		// 正常系
		{"male", "male", "male", nil},
		{"female", "female", "female", nil},
		{"other", "other", "other", nil},
		// 異常系
		{"empty", "", "", domainErrors.ErrInvalidGender},
		{"invalid", "unknown", "", domainErrors.ErrInvalidGender},
		{"uppercase MALE", "MALE", "", domainErrors.ErrInvalidGender},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vo.NewGender(tt.input)

			if err != tt.wantErr {
				t.Errorf("NewGender(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if err == nil && got.String() != tt.wantGender {
				t.Errorf("NewGender(%q).String() = %v, want %v", tt.input, got.String(), tt.wantGender)
			}
		})
	}
}
