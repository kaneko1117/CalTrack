package vo_test

import (
	"strings"
	"testing"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
)

func TestNewEmail(t *testing.T) {
	// 境界値テスト用の長いメール作成
	local64 := strings.Repeat("a", 64)
	domain185 := strings.Repeat("b", 185) + ".com" // 254 - 64 - 1(@) - 4(.com) = 185
	email254 := local64 + "@" + domain185
	domain186 := strings.Repeat("b", 186) + ".com"
	email255 := local64 + "@" + domain186

	tests := []struct {
		name      string
		input     string
		wantEmail string
		wantErr   error
	}{
		// 正常系
		{"valid email", "user@example.com", "user@example.com", nil},
		{"valid email with subdomain", "user@mail.example.com", "user@mail.example.com", nil},
		{"valid email with plus", "user+tag@example.com", "user+tag@example.com", nil},
		// 異常系
		{"empty string", "", "", domainErrors.ErrEmailRequired},
		{"no at sign", "userexample.com", "", domainErrors.ErrInvalidEmailFormat},
		{"no domain", "user@", "", domainErrors.ErrInvalidEmailFormat},
		{"no local part", "@example.com", "", domainErrors.ErrInvalidEmailFormat},
		// 境界値
		{"max length 254", email254, email254, nil},
		{"exceeds 254", email255, "", domainErrors.ErrEmailTooLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vo.NewEmail(tt.input)

			if err != tt.wantErr {
				t.Errorf("NewEmail(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if err == nil && got.String() != tt.wantEmail {
				t.Errorf("NewEmail(%q).String() = %v, want %v", tt.input, got.String(), tt.wantEmail)
			}
		})
	}
}

func TestEmail_Equals(t *testing.T) {
	tests := []struct {
		name   string
		email1 string
		email2 string
		want   bool
	}{
		{"same value", "user@example.com", "user@example.com", true},
		{"different value", "user1@example.com", "user2@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e1, _ := vo.NewEmail(tt.email1)
			e2, _ := vo.NewEmail(tt.email2)

			if got := e1.Equals(e2); got != tt.want {
				t.Errorf("Email.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}
