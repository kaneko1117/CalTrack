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
		{"maleは有効", "male", "male", nil},
		{"femaleは有効", "female", "female", nil},
		{"otherは有効", "other", "other", nil},
		// 異常系
		{"空文字はエラー", "", "", domainErrors.ErrInvalidGender},
		{"無効な値はエラー", "unknown", "", domainErrors.ErrInvalidGender},
		{"大文字MALEはエラー", "MALE", "", domainErrors.ErrInvalidGender},
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
