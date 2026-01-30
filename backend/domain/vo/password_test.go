package vo_test

import (
	"testing"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
)

func TestNewPassword(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		// 正常系
		{"valid 8 chars", "12345678", nil},
		{"valid long password", "verylongpassword123", nil},
		// 異常系
		{"empty string", "", domainErrors.ErrPasswordRequired},
		{"too short 7 chars", "1234567", domainErrors.ErrPasswordTooShort},
		// 境界値
		{"exactly 8 chars", "abcdefgh", nil},
		{"only 7 chars", "abcdefg", domainErrors.ErrPasswordTooShort},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := vo.NewPassword(tt.input)

			if err != tt.wantErr {
				t.Errorf("NewPassword(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestPassword_Hash(t *testing.T) {
	password, _ := vo.NewPassword("12345678")

	hashed, err := password.Hash()

	if err != nil {
		t.Errorf("Hash() error = %v", err)
	}
	if hashed.String() == "" {
		t.Error("Hash() should return non-empty hash")
	}
	if hashed.String() == "12345678" {
		t.Error("Hash() should not equal plain password")
	}
}

func TestHashedPassword_Compare(t *testing.T) {
	password, _ := vo.NewPassword("12345678")
	hashed, _ := password.Hash()

	tests := []struct {
		name     string
		input    string
		wantMatch bool
	}{
		{"correct password", "12345678", true},
		{"wrong password", "wrongpassword", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, _ := vo.NewPassword(tt.input)
			if got := hashed.Compare(p); got != tt.wantMatch {
				t.Errorf("Compare() = %v, want %v", got, tt.wantMatch)
			}
		})
	}
}

func TestNewHashedPassword_Reconstruct(t *testing.T) {
	password, _ := vo.NewPassword("12345678")
	original, _ := password.Hash()

	reconstructed := vo.NewHashedPassword(original.String())

	samePassword, _ := vo.NewPassword("12345678")
	if !reconstructed.Compare(samePassword) {
		t.Error("Reconstructed hash should match original password")
	}
}
