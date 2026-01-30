package vo_test

import (
	"strings"
	"testing"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
)

func TestNewNickname(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantNickname string
		wantErr      error
	}{
		// 正常系
		{"valid nickname", "John", "John", nil},
		{"valid with spaces", "John Doe", "John Doe", nil},
		// 異常系
		{"empty string", "", "", domainErrors.ErrNicknameRequired},
		{"too long 51 chars", strings.Repeat("a", 51), "", domainErrors.ErrNicknameTooLong},
		// 境界値
		{"min length 1 char", "A", "A", nil},
		{"max length 50 chars", strings.Repeat("a", 50), strings.Repeat("a", 50), nil},
		{"exceeds 50 chars", strings.Repeat("a", 51), "", domainErrors.ErrNicknameTooLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vo.NewNickname(tt.input)

			if err != tt.wantErr {
				t.Errorf("NewNickname(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if err == nil && got.String() != tt.wantNickname {
				t.Errorf("NewNickname(%q).String() = %v, want %v", tt.input, got.String(), tt.wantNickname)
			}
		})
	}
}
