package vo_test

import (
	"testing"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
)

func TestNewPassword_Valid8Chars_ReturnsPassword(t *testing.T) {
	_, err := vo.NewPassword("12345678")

	if err != nil {
		t.Errorf("NewPassword should not return error for 8 char password, got: %v", err)
	}
}

func TestNewPassword_ValidLongPassword_ReturnsPassword(t *testing.T) {
	_, err := vo.NewPassword("verylongpassword123")

	if err != nil {
		t.Errorf("NewPassword should not return error for long password, got: %v", err)
	}
}

func TestNewPassword_EmptyString_ReturnsError(t *testing.T) {
	_, err := vo.NewPassword("")

	if err == nil {
		t.Error("NewPassword should return error for empty string")
	}
	if err != domainErrors.ErrPasswordRequired {
		t.Errorf("NewPassword should return ErrPasswordRequired, got: %v", err)
	}
}

func TestNewPassword_TooShort7Chars_ReturnsError(t *testing.T) {
	_, err := vo.NewPassword("1234567")

	if err == nil {
		t.Error("NewPassword should return error for 7 char password")
	}
	if err != domainErrors.ErrPasswordTooShort {
		t.Errorf("NewPassword should return ErrPasswordTooShort, got: %v", err)
	}
}

func TestNewPassword_Exactly8Chars_ReturnsPassword(t *testing.T) {
	_, err := vo.NewPassword("12345678")

	if err != nil {
		t.Errorf("NewPassword should accept exactly 8 chars, got: %v", err)
	}
}

func TestNewPassword_Only7Chars_ReturnsError(t *testing.T) {
	_, err := vo.NewPassword("1234567")

	if err == nil {
		t.Error("NewPassword should reject 7 char password")
	}
}

func TestPassword_Hash_ReturnsHashedPassword(t *testing.T) {
	password, _ := vo.NewPassword("12345678")

	hashed, err := password.Hash()

	if err != nil {
		t.Errorf("Hash should not return error, got: %v", err)
	}
	if hashed.String() == "" {
		t.Error("HashedPassword should not be empty")
	}
	if hashed.String() == "12345678" {
		t.Error("HashedPassword should not equal plain password")
	}
}

func TestHashedPassword_Compare_CorrectPassword_ReturnsTrue(t *testing.T) {
	password, _ := vo.NewPassword("12345678")
	hashed, _ := password.Hash()

	samePassword, _ := vo.NewPassword("12345678")
	result := hashed.Compare(samePassword)

	if !result {
		t.Error("Compare should return true for correct password")
	}
}

func TestHashedPassword_Compare_WrongPassword_ReturnsFalse(t *testing.T) {
	password, _ := vo.NewPassword("12345678")
	hashed, _ := password.Hash()

	wrongPassword, _ := vo.NewPassword("wrongpassword")
	result := hashed.Compare(wrongPassword)

	if result {
		t.Error("Compare should return false for wrong password")
	}
}

func TestNewHashedPassword_ReconstructsFromHash(t *testing.T) {
	password, _ := vo.NewPassword("12345678")
	original, _ := password.Hash()

	reconstructed := vo.NewHashedPassword(original.String())

	samePassword, _ := vo.NewPassword("12345678")
	if !reconstructed.Compare(samePassword) {
		t.Error("Reconstructed HashedPassword should compare correctly")
	}
}
